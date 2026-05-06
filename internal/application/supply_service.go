package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// SupplyService handles inbound supply workflows.
type SupplyService struct {
	provider      RepositoryProvider
	transactions  TransactionManager
	notifications NotificationPublisher
}

// NewSupplyService creates a supply service.
func NewSupplyService(provider RepositoryProvider, transactions TransactionManager, publishers ...NotificationPublisher) *SupplyService {
	return &SupplyService{
		provider:      provider,
		transactions:  transactions,
		notifications: notificationPublisherFrom(publishers...),
	}
}

// Create creates a draft supply for an active warehouse.
func (service *SupplyService) Create(ctx context.Context, input CreateSupplyInput) (*domain.Supply, error) {
	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, input.CreatedByUserID); err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, input.WarehouseID); err != nil {
		return nil, err
	}

	supply := domain.Supply{
		WarehouseID:     input.WarehouseID,
		Status:          domain.SupplyStatusDraft,
		CreatedByUserID: input.CreatedByUserID,
	}
	if err := repos.Supplies.Create(ctx, &supply); err != nil {
		return nil, err
	}

	return &supply, nil
}

// AddItem adds a product row to an editable supply.
func (service *SupplyService) AddItem(ctx context.Context, supplyID domain.ID, input SupplyItemInput) (*domain.SupplyItem, error) {
	if err := requireID(supplyID, "supply id"); err != nil {
		return nil, err
	}
	if err := validateSupplyItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	supply, err := service.requireEditableSupply(ctx, repos, supplyID)
	if err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, supply.WarehouseID); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureSupplyProductNotDuplicated(ctx, repos.Supplies, supplyID, 0, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.SupplyItem{
		SupplyID:       supplyID,
		ProductID:      input.ProductID,
		Quantity:       input.Quantity,
		UnitPriceCents: input.UnitPriceCents,
	}
	if err := repos.Supplies.AddItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItem updates a product row on an editable supply.
func (service *SupplyService) UpdateItem(ctx context.Context, supplyID domain.ID, itemID domain.ID, input SupplyItemInput) (*domain.SupplyItem, error) {
	if err := requireID(supplyID, "supply id"); err != nil {
		return nil, err
	}
	if err := requireID(itemID, "supply item id"); err != nil {
		return nil, err
	}
	if err := validateSupplyItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	supply, err := service.requireEditableSupply(ctx, repos, supplyID)
	if err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, supply.WarehouseID); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureSupplyProductNotDuplicated(ctx, repos.Supplies, supplyID, itemID, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.SupplyItem{
		ID:             itemID,
		SupplyID:       supplyID,
		ProductID:      input.ProductID,
		Quantity:       input.Quantity,
		UnitPriceCents: input.UnitPriceCents,
	}
	if err := repos.Supplies.UpdateItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem removes a product row from an editable supply.
func (service *SupplyService) DeleteItem(ctx context.Context, supplyID domain.ID, itemID domain.ID) error {
	if err := requireID(supplyID, "supply id"); err != nil {
		return err
	}
	if err := requireID(itemID, "supply item id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableSupply(ctx, repos, supplyID); err != nil {
		return err
	}
	if err := ensureSupplyItemBelongs(ctx, repos.Supplies, supplyID, itemID); err != nil {
		return err
	}

	return repos.Supplies.DeleteItem(ctx, itemID)
}

// Confirm confirms a supply and increases warehouse stock inside one transaction.
func (service *SupplyService) Confirm(ctx context.Context, id domain.ID, confirmedByUserID domain.ID) (*domain.Supply, error) {
	if err := requireID(id, "supply id"); err != nil {
		return nil, err
	}
	if service.transactions == nil {
		return nil, businessRuleError(ErrBusinessRule, "transaction manager is required")
	}

	var confirmed *domain.Supply
	var stockEvents []NotificationEvent
	err := service.transactions.WithinTransaction(ctx, func(txCtx context.Context, repos Repositories) error {
		if _, err := requireActiveUser(txCtx, repos.Users, confirmedByUserID); err != nil {
			return err
		}
		supply, err := repos.Supplies.GetByID(txCtx, id)
		if err != nil {
			return err
		}
		if !supply.Status.CanConfirm() {
			return businessRuleError(ErrOperationLocked, "completed or cancelled supply cannot be confirmed")
		}
		if !supply.HasItems() {
			return businessRuleError(ErrBusinessRule, "cannot confirm supply without items")
		}
		if !supply.HasValidItems() {
			return validationError("supply items must have positive quantities and non-negative prices")
		}
		if _, err := requireActiveWarehouse(txCtx, repos.Warehouses, supply.WarehouseID); err != nil {
			return err
		}

		for _, item := range supply.Items {
			product, err := requireActiveProduct(txCtx, repos.Products, item.ProductID)
			if err != nil {
				return err
			}
			stockItem, err := repos.Stock.Increment(txCtx, supply.WarehouseID, item.ProductID, item.Quantity)
			if err != nil {
				return err
			}
			stockEvents = appendStockStateEvent(stockEvents, *stockItem, *product, NotificationSupplyConfirmed)
		}

		now := nowUTC()
		supply.Status = domain.SupplyStatusCompleted
		supply.ConfirmedByUserID = &confirmedByUserID
		supply.ConfirmedAt = &now
		if err := repos.Supplies.Update(txCtx, supply); err != nil {
			return err
		}

		confirmed, err = repos.Supplies.GetByID(txCtx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	service.notifications.Publish(ctx, supplyConfirmedEvent(*confirmed))
	publishAll(ctx, service.notifications, stockEvents)

	return confirmed, nil
}

// Cancel cancels a draft supply without changing stock.
func (service *SupplyService) Cancel(ctx context.Context, id domain.ID, cancelledByUserID domain.ID) (*domain.Supply, error) {
	if err := requireID(id, "supply id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, cancelledByUserID); err != nil {
		return nil, err
	}
	supply, err := repos.Supplies.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !supply.Status.CanCancel() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled supply cannot be cancelled")
	}

	now := nowUTC()
	supply.Status = domain.SupplyStatusCancelled
	supply.CancelledByUserID = &cancelledByUserID
	supply.CancelledAt = &now
	if err := repos.Supplies.Update(ctx, supply); err != nil {
		return nil, err
	}

	return repos.Supplies.GetByID(ctx, id)
}

// Delete removes a draft supply.
func (service *SupplyService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "supply id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableSupply(ctx, repos, id); err != nil {
		return err
	}

	return repos.Supplies.Delete(ctx, id)
}

// GetByID returns a supply by identifier.
func (service *SupplyService) GetByID(ctx context.Context, id domain.ID) (*domain.Supply, error) {
	if err := requireID(id, "supply id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Supplies.GetByID(ctx, id)
}

// List returns supplies matching the filter.
func (service *SupplyService) List(ctx context.Context, filter SupplyFilter) ([]domain.Supply, error) {
	return repositories(service.provider).Supplies.List(ctx, filter)
}

func (service *SupplyService) requireEditableSupply(ctx context.Context, repos Repositories, id domain.ID) (*domain.Supply, error) {
	supply, err := repos.Supplies.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !supply.CanBeEdited() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled supply cannot be edited")
	}

	return supply, nil
}

func validateSupplyItemInput(input SupplyItemInput) error {
	if err := requireID(input.ProductID, "product id"); err != nil {
		return err
	}
	if err := ensurePositiveQuantity(input.Quantity); err != nil {
		return err
	}
	return ensureNonNegativePrice(input.UnitPriceCents)
}

func ensureSupplyProductNotDuplicated(ctx context.Context, repository SupplyRepository, supplyID domain.ID, currentItemID domain.ID, productID domain.ID) error {
	items, err := repository.ListItems(ctx, supplyID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ProductID == productID && item.ID != currentItemID {
			return conflictError(ErrConflict, "product already exists in supply")
		}
	}

	return nil
}

func ensureSupplyItemBelongs(ctx context.Context, repository SupplyRepository, supplyID domain.ID, itemID domain.ID) error {
	items, err := repository.ListItems(ctx, supplyID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == itemID {
			return nil
		}
	}

	return ErrNotFound
}
