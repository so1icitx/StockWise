package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

type unitFixture struct {
	store             *fakeStore
	publisher         *recordingPublisher
	services          application.Services
	admin             domain.User
	manager           domain.User
	operator          domain.User
	warehouse         domain.Warehouse
	targetWarehouse   domain.Warehouse
	inactiveWarehouse domain.Warehouse
	category          domain.Category
	product           domain.Product
	secondaryProduct  domain.Product
}

func TestProductCreateRejectsDuplicateSKU(t *testing.T) {
	fixture := newUnitFixture(t)

	_, err := fixture.services.Products.Create(context.Background(), application.CreateProductInput{
		Name:              "Duplicate Scanner",
		SKU:               fixture.product.SKU,
		CategoryID:        fixture.category.ID,
		UnitOfMeasure:     "pcs",
		MinStockThreshold: 3,
	})

	assertErrorIs(t, err, application.ErrDuplicateSKU)
}

func TestSupplyConfirmRejectsEmptySupply(t *testing.T) {
	fixture := newUnitFixture(t)
	supply := fixture.createSupply(t)

	_, err := fixture.services.Supplies.Confirm(context.Background(), supply.ID, fixture.manager.ID)

	assertErrorIs(t, err, application.ErrBusinessRule)
	if events := fixture.publisher.eventsNamed(application.NotificationSupplyConfirmed); len(events) != 0 {
		t.Fatalf("expected no supply.confirmed events, got %d", len(events))
	}
}

func TestSupplyAddItemRejectsNegativePrice(t *testing.T) {
	fixture := newUnitFixture(t)
	supply := fixture.createSupply(t)

	_, err := fixture.services.Supplies.AddItem(context.Background(), supply.ID, application.SupplyItemInput{
		ProductID:      fixture.product.ID,
		Quantity:       1,
		UnitPriceCents: -1,
	})

	assertErrorIs(t, err, application.ErrValidation)
}

func TestSupplyConfirmIncreasesStock(t *testing.T) {
	fixture := newUnitFixture(t)
	supply := fixture.createSupply(t)
	fixture.addSupplyItem(t, supply.ID, fixture.product.ID, 7, 119900)

	confirmed, err := fixture.services.Supplies.Confirm(context.Background(), supply.ID, fixture.manager.ID)
	if err != nil {
		t.Fatalf("confirm supply: %v", err)
	}

	if confirmed.Status != domain.SupplyStatusCompleted {
		t.Fatalf("expected completed supply, got %s", confirmed.Status)
	}
	assertStockQuantity(t, fixture, fixture.warehouse.ID, fixture.product.ID, 17)
	events := fixture.publisher.eventsNamed(application.NotificationSupplyConfirmed)
	if len(events) != 1 {
		t.Fatalf("expected one supply.confirmed event, got %d", len(events))
	}
	if got := events[0].Data["supply_id"]; got != supply.ID {
		t.Fatalf("expected supply_id %d in notification, got %v", supply.ID, got)
	}
}

func TestOutboundExecuteRejectsInsufficientStock(t *testing.T) {
	fixture := newUnitFixture(t)
	request := fixture.createApprovedOutbound(t, 99)

	_, err := fixture.services.OutboundRequests.Execute(context.Background(), request.ID, fixture.operator.ID)

	assertErrorIs(t, err, application.ErrInsufficientStock)
	assertStockQuantity(t, fixture, fixture.warehouse.ID, fixture.product.ID, 10)
	if events := fixture.publisher.eventsNamed(application.NotificationOutboundCompleted); len(events) != 0 {
		t.Fatalf("expected no outbound.completed events, got %d", len(events))
	}
}

func TestOutboundExecuteDecreasesStock(t *testing.T) {
	fixture := newUnitFixture(t)
	request := fixture.createApprovedOutbound(t, 4)

	completed, err := fixture.services.OutboundRequests.Execute(context.Background(), request.ID, fixture.operator.ID)
	if err != nil {
		t.Fatalf("execute outbound request: %v", err)
	}

	if completed.Status != domain.OutboundRequestStatusCompleted {
		t.Fatalf("expected completed outbound request, got %s", completed.Status)
	}
	assertStockQuantity(t, fixture, fixture.warehouse.ID, fixture.product.ID, 6)
	events := fixture.publisher.eventsNamed(application.NotificationOutboundCompleted)
	if len(events) != 1 {
		t.Fatalf("expected one outbound.completed event, got %d", len(events))
	}
}

func TestTransferCreateRejectsSameWarehouse(t *testing.T) {
	fixture := newUnitFixture(t)

	_, err := fixture.services.Transfers.Create(context.Background(), application.CreateTransferInput{
		SourceWarehouseID: fixture.warehouse.ID,
		TargetWarehouseID: fixture.warehouse.ID,
		CreatedByUserID:   fixture.operator.ID,
	})

	assertErrorIs(t, err, application.ErrBusinessRule)
}

func TestTransferConfirmRejectsEmptyTransfer(t *testing.T) {
	fixture := newUnitFixture(t)
	transfer := fixture.createTransfer(t)

	_, err := fixture.services.Transfers.Confirm(context.Background(), transfer.ID, fixture.manager.ID)

	assertErrorIs(t, err, application.ErrBusinessRule)
}

func TestTransferConfirmMovesStockBetweenWarehouses(t *testing.T) {
	fixture := newUnitFixture(t)
	transfer := fixture.createTransfer(t)
	fixture.addTransferItem(t, transfer.ID, fixture.product.ID, 3)

	completed, err := fixture.services.Transfers.Confirm(context.Background(), transfer.ID, fixture.manager.ID)
	if err != nil {
		t.Fatalf("confirm transfer: %v", err)
	}

	if completed.Status != domain.TransferStatusCompleted {
		t.Fatalf("expected completed transfer, got %s", completed.Status)
	}
	assertStockQuantity(t, fixture, fixture.warehouse.ID, fixture.product.ID, 7)
	assertStockQuantity(t, fixture, fixture.targetWarehouse.ID, fixture.product.ID, 3)
	events := fixture.publisher.eventsNamed(application.NotificationTransferConfirmed)
	if len(events) != 1 {
		t.Fatalf("expected one transfer.confirmed event, got %d", len(events))
	}
}

func TestWarehouseDeleteRejectsWarehouseWithStock(t *testing.T) {
	fixture := newUnitFixture(t)

	err := fixture.services.Warehouses.Delete(context.Background(), fixture.warehouse.ID)

	assertErrorIs(t, err, application.ErrDeleteRestricted)
	if _, err := fixture.services.Warehouses.GetByID(context.Background(), fixture.warehouse.ID); err != nil {
		t.Fatalf("warehouse should still exist after rejected delete: %v", err)
	}
}

func TestInactiveWarehouseBlocksOperations(t *testing.T) {
	fixture := newUnitFixture(t)

	_, err := fixture.services.Supplies.Create(context.Background(), application.CreateSupplyInput{
		WarehouseID:     fixture.inactiveWarehouse.ID,
		CreatedByUserID: fixture.operator.ID,
	})
	assertErrorIs(t, err, application.ErrInactiveWarehouse)

	_, err = fixture.services.OutboundRequests.Create(context.Background(), application.CreateOutboundRequestInput{
		WarehouseID:     fixture.inactiveWarehouse.ID,
		CreatedByUserID: fixture.operator.ID,
	})
	assertErrorIs(t, err, application.ErrInactiveWarehouse)

	_, err = fixture.services.Transfers.Create(context.Background(), application.CreateTransferInput{
		SourceWarehouseID: fixture.warehouse.ID,
		TargetWarehouseID: fixture.inactiveWarehouse.ID,
		CreatedByUserID:   fixture.operator.ID,
	})
	assertErrorIs(t, err, application.ErrInactiveWarehouse)
}

func TestCompletedSupplyCannotBeEdited(t *testing.T) {
	fixture := newUnitFixture(t)
	supply := fixture.createSupply(t)
	fixture.addSupplyItem(t, supply.ID, fixture.product.ID, 1, 119900)
	if _, err := fixture.services.Supplies.Confirm(context.Background(), supply.ID, fixture.manager.ID); err != nil {
		t.Fatalf("confirm supply: %v", err)
	}

	_, err := fixture.services.Supplies.AddItem(context.Background(), supply.ID, application.SupplyItemInput{
		ProductID:      fixture.secondaryProduct.ID,
		Quantity:       1,
		UnitPriceCents: 500,
	})

	assertErrorIs(t, err, application.ErrOperationLocked)
}

func TestCompletedOutboundRequestCannotBeEdited(t *testing.T) {
	fixture := newUnitFixture(t)
	request := fixture.createApprovedOutbound(t, 1)
	if _, err := fixture.services.OutboundRequests.Execute(context.Background(), request.ID, fixture.operator.ID); err != nil {
		t.Fatalf("execute outbound request: %v", err)
	}

	_, err := fixture.services.OutboundRequests.AddItem(context.Background(), request.ID, application.OutboundRequestItemInput{
		ProductID: fixture.secondaryProduct.ID,
		Quantity:  1,
	})

	assertErrorIs(t, err, application.ErrOperationLocked)
}

func TestCompletedTransferCannotBeEdited(t *testing.T) {
	fixture := newUnitFixture(t)
	transfer := fixture.createTransfer(t)
	fixture.addTransferItem(t, transfer.ID, fixture.product.ID, 1)
	if _, err := fixture.services.Transfers.Confirm(context.Background(), transfer.ID, fixture.manager.ID); err != nil {
		t.Fatalf("confirm transfer: %v", err)
	}

	_, err := fixture.services.Transfers.AddItem(context.Background(), transfer.ID, application.TransferItemInput{
		ProductID: fixture.secondaryProduct.ID,
		Quantity:  1,
	})

	assertErrorIs(t, err, application.ErrOperationLocked)
}

func TestOutboundExecutionPublishesLowAndOutOfStockNotifications(t *testing.T) {
	fixture := newUnitFixture(t)
	fixture.setStock(fixture.warehouse.ID, fixture.product.ID, 6)

	lowRequest := fixture.createApprovedOutbound(t, 2)
	if _, err := fixture.services.OutboundRequests.Execute(context.Background(), lowRequest.ID, fixture.operator.ID); err != nil {
		t.Fatalf("execute low-stock outbound request: %v", err)
	}
	lowEvents := fixture.publisher.eventsNamed(application.NotificationStockLow)
	if len(lowEvents) != 1 {
		t.Fatalf("expected one stock.low event, got %d", len(lowEvents))
	}
	if got := lowEvents[0].Data["quantity"]; got != int64(4) {
		t.Fatalf("expected low-stock quantity 4, got %v", got)
	}

	outRequest := fixture.createApprovedOutbound(t, 4)
	if _, err := fixture.services.OutboundRequests.Execute(context.Background(), outRequest.ID, fixture.operator.ID); err != nil {
		t.Fatalf("execute out-of-stock outbound request: %v", err)
	}
	outEvents := fixture.publisher.eventsNamed(application.NotificationStockOut)
	if len(outEvents) != 1 {
		t.Fatalf("expected one stock.out event, got %d", len(outEvents))
	}
	if got := outEvents[0].Data["quantity"]; got != int64(0) {
		t.Fatalf("expected out-of-stock quantity 0, got %v", got)
	}
}

func TestStockServiceClassifiesLowAndOutOfStockRows(t *testing.T) {
	fixture := newUnitFixture(t)
	fixture.setStock(fixture.warehouse.ID, fixture.product.ID, 5)
	fixture.setStock(fixture.targetWarehouse.ID, fixture.secondaryProduct.ID, 0)

	statuses, err := fixture.services.Stock.GetLowStock(context.Background())
	if err != nil {
		t.Fatalf("get low stock: %v", err)
	}

	states := map[domain.StockState]bool{}
	for _, status := range statuses {
		states[status.State] = true
	}
	if !states[domain.StockStateLow] {
		t.Fatal("expected a low-stock row")
	}
	if !states[domain.StockStateOut] {
		t.Fatal("expected an out-of-stock row")
	}
}

func TestWarehouseDeactivatePublishesNotification(t *testing.T) {
	fixture := newUnitFixture(t)

	if err := fixture.services.Warehouses.Deactivate(context.Background(), fixture.targetWarehouse.ID); err != nil {
		t.Fatalf("deactivate warehouse: %v", err)
	}

	events := fixture.publisher.eventsNamed(application.NotificationWarehouseDeactivated)
	if len(events) != 1 {
		t.Fatalf("expected one warehouse.deactivated event, got %d", len(events))
	}
	if got := events[0].Data["warehouse_id"]; got != fixture.targetWarehouse.ID {
		t.Fatalf("expected warehouse_id %d, got %v", fixture.targetWarehouse.ID, got)
	}
}

func newUnitFixture(t *testing.T) *unitFixture {
	t.Helper()

	store := newFakeStore()
	publisher := &recordingPublisher{}
	fixture := &unitFixture{
		store:     store,
		publisher: publisher,
	}

	fixture.admin = store.addUser(domain.User{Name: "Admin User", Email: "admin@stockwise.test", Role: domain.UserRoleAdmin, IsActive: true})
	fixture.manager = store.addUser(domain.User{Name: "Manager User", Email: "manager@stockwise.test", Role: domain.UserRoleManager, IsActive: true})
	fixture.operator = store.addUser(domain.User{Name: "Operator User", Email: "operator@stockwise.test", Role: domain.UserRoleOperator, IsActive: true})
	fixture.warehouse = store.addWarehouse(domain.Warehouse{Name: "Main Warehouse", Code: "MAIN", Location: "Sofia", IsActive: true})
	fixture.targetWarehouse = store.addWarehouse(domain.Warehouse{Name: "Secondary Warehouse", Code: "SECOND", Location: "Plovdiv", IsActive: true})
	fixture.inactiveWarehouse = store.addWarehouse(domain.Warehouse{Name: "Closed Warehouse", Code: "CLOSED", Location: "Varna", IsActive: false})
	fixture.category = store.addCategory(domain.Category{Name: "Electronics", Description: "Electronic devices", IsActive: true})
	fixture.product = store.addProduct(domain.Product{
		Name:              "Wireless Scanner",
		SKU:               "ELEC-001",
		CategoryID:        fixture.category.ID,
		UnitOfMeasure:     "pcs",
		MinStockThreshold: 5,
		IsActive:          true,
	})
	fixture.secondaryProduct = store.addProduct(domain.Product{
		Name:              "Barcode Cable",
		SKU:               "ELEC-002",
		CategoryID:        fixture.category.ID,
		UnitOfMeasure:     "pcs",
		MinStockThreshold: 3,
		IsActive:          true,
	})
	store.addStock(domain.StockItem{
		WarehouseID: fixture.warehouse.ID,
		ProductID:   fixture.product.ID,
		Quantity:    10,
	})

	fixture.services = application.NewServices(store, store, publisher)
	return fixture
}

func (fixture *unitFixture) createSupply(t *testing.T) *domain.Supply {
	t.Helper()
	supply, err := fixture.services.Supplies.Create(context.Background(), application.CreateSupplyInput{
		WarehouseID:     fixture.warehouse.ID,
		CreatedByUserID: fixture.operator.ID,
	})
	if err != nil {
		t.Fatalf("create supply: %v", err)
	}
	return supply
}

func (fixture *unitFixture) addSupplyItem(t *testing.T, supplyID domain.ID, productID domain.ID, quantity int64, unitPriceCents int64) *domain.SupplyItem {
	t.Helper()
	item, err := fixture.services.Supplies.AddItem(context.Background(), supplyID, application.SupplyItemInput{
		ProductID:      productID,
		Quantity:       quantity,
		UnitPriceCents: unitPriceCents,
	})
	if err != nil {
		t.Fatalf("add supply item: %v", err)
	}
	return item
}

func (fixture *unitFixture) createApprovedOutbound(t *testing.T, quantity int64) *domain.OutboundRequest {
	t.Helper()
	request, err := fixture.services.OutboundRequests.Create(context.Background(), application.CreateOutboundRequestInput{
		WarehouseID:     fixture.warehouse.ID,
		CreatedByUserID: fixture.operator.ID,
	})
	if err != nil {
		t.Fatalf("create outbound request: %v", err)
	}
	if _, err := fixture.services.OutboundRequests.AddItem(context.Background(), request.ID, application.OutboundRequestItemInput{
		ProductID: fixture.product.ID,
		Quantity:  quantity,
	}); err != nil {
		t.Fatalf("add outbound request item: %v", err)
	}
	approved, err := fixture.services.OutboundRequests.Approve(context.Background(), request.ID, fixture.manager.ID)
	if err != nil {
		t.Fatalf("approve outbound request: %v", err)
	}
	return approved
}

func (fixture *unitFixture) createTransfer(t *testing.T) *domain.Transfer {
	t.Helper()
	transfer, err := fixture.services.Transfers.Create(context.Background(), application.CreateTransferInput{
		SourceWarehouseID: fixture.warehouse.ID,
		TargetWarehouseID: fixture.targetWarehouse.ID,
		CreatedByUserID:   fixture.operator.ID,
	})
	if err != nil {
		t.Fatalf("create transfer: %v", err)
	}
	return transfer
}

func (fixture *unitFixture) addTransferItem(t *testing.T, transferID domain.ID, productID domain.ID, quantity int64) *domain.TransferItem {
	t.Helper()
	item, err := fixture.services.Transfers.AddItem(context.Background(), transferID, application.TransferItemInput{
		ProductID: productID,
		Quantity:  quantity,
	})
	if err != nil {
		t.Fatalf("add transfer item: %v", err)
	}
	return item
}

func (fixture *unitFixture) setStock(warehouseID domain.ID, productID domain.ID, quantity int64) {
	fixture.store.addStock(domain.StockItem{
		WarehouseID: warehouseID,
		ProductID:   productID,
		Quantity:    quantity,
	})
}

func assertStockQuantity(t *testing.T, fixture *unitFixture, warehouseID domain.ID, productID domain.ID, want int64) {
	t.Helper()
	stockItem, err := fixture.services.Stock.GetByWarehouseAndProduct(context.Background(), warehouseID, productID)
	if err != nil {
		t.Fatalf("get stock item: %v", err)
	}
	if stockItem.Quantity != want {
		t.Fatalf("expected stock quantity %d, got %d", want, stockItem.Quantity)
	}
}

func assertErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}
