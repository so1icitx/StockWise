package unit

import (
	"context"
	"strings"
	"time"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

type stockKey struct {
	warehouseID domain.ID
	productID   domain.ID
}

type fakeStore struct {
	nextID        domain.ID
	users         map[domain.ID]domain.User
	warehouses    map[domain.ID]domain.Warehouse
	categories    map[domain.ID]domain.Category
	products      map[domain.ID]domain.Product
	stock         map[stockKey]domain.StockItem
	supplies      map[domain.ID]domain.Supply
	supplyItems   map[domain.ID]domain.SupplyItem
	outbounds     map[domain.ID]domain.OutboundRequest
	outboundItems map[domain.ID]domain.OutboundRequestItem
	transfers     map[domain.ID]domain.Transfer
	transferItems map[domain.ID]domain.TransferItem
}

type recordingPublisher struct {
	events []application.NotificationEvent
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		nextID:        1,
		users:         map[domain.ID]domain.User{},
		warehouses:    map[domain.ID]domain.Warehouse{},
		categories:    map[domain.ID]domain.Category{},
		products:      map[domain.ID]domain.Product{},
		stock:         map[stockKey]domain.StockItem{},
		supplies:      map[domain.ID]domain.Supply{},
		supplyItems:   map[domain.ID]domain.SupplyItem{},
		outbounds:     map[domain.ID]domain.OutboundRequest{},
		outboundItems: map[domain.ID]domain.OutboundRequestItem{},
		transfers:     map[domain.ID]domain.Transfer{},
		transferItems: map[domain.ID]domain.TransferItem{},
	}
}

func (store *fakeStore) Repositories() application.Repositories {
	return application.Repositories{
		Users:            fakeUserRepository{store: store},
		Warehouses:       fakeWarehouseRepository{store: store},
		Categories:       fakeCategoryRepository{store: store},
		Products:         fakeProductRepository{store: store},
		Stock:            fakeStockRepository{store: store},
		Supplies:         fakeSupplyRepository{store: store},
		OutboundRequests: fakeOutboundRepository{store: store},
		Transfers:        fakeTransferRepository{store: store},
		Movements:        fakeMovementRepository{store: store},
	}
}

func (store *fakeStore) WithinTransaction(ctx context.Context, fn func(context.Context, application.Repositories) error) error {
	return fn(ctx, store.Repositories())
}

func (store *fakeStore) assignID(id *domain.ID) {
	if !id.IsZero() {
		if *id >= store.nextID {
			store.nextID = *id + 1
		}
		return
	}

	*id = store.nextID
	store.nextID++
}

func (store *fakeStore) addUser(user domain.User) domain.User {
	store.assignID(&user.ID)
	user.CreatedAt = defaultTime(user.CreatedAt)
	user.UpdatedAt = defaultTime(user.UpdatedAt)
	store.users[user.ID] = user
	return user
}

func (store *fakeStore) addWarehouse(warehouse domain.Warehouse) domain.Warehouse {
	store.assignID(&warehouse.ID)
	warehouse.CreatedAt = defaultTime(warehouse.CreatedAt)
	warehouse.UpdatedAt = defaultTime(warehouse.UpdatedAt)
	store.warehouses[warehouse.ID] = warehouse
	return warehouse
}

func (store *fakeStore) addCategory(category domain.Category) domain.Category {
	store.assignID(&category.ID)
	category.CreatedAt = defaultTime(category.CreatedAt)
	category.UpdatedAt = defaultTime(category.UpdatedAt)
	store.categories[category.ID] = category
	return category
}

func (store *fakeStore) addProduct(product domain.Product) domain.Product {
	store.assignID(&product.ID)
	product.CreatedAt = defaultTime(product.CreatedAt)
	product.UpdatedAt = defaultTime(product.UpdatedAt)
	store.products[product.ID] = product
	return product
}

func (store *fakeStore) addStock(stockItem domain.StockItem) domain.StockItem {
	store.assignID(&stockItem.ID)
	stockItem.CreatedAt = defaultTime(stockItem.CreatedAt)
	stockItem.UpdatedAt = defaultTime(stockItem.UpdatedAt)
	store.stock[stockKey{warehouseID: stockItem.WarehouseID, productID: stockItem.ProductID}] = stockItem
	return stockItem
}

func defaultTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Date(2026, 4, 30, 9, 0, 0, 0, time.UTC)
	}

	return value
}

func (publisher *recordingPublisher) Publish(_ context.Context, event application.NotificationEvent) {
	publisher.events = append(publisher.events, event)
}

func (publisher *recordingPublisher) eventsNamed(name application.NotificationEventName) []application.NotificationEvent {
	events := make([]application.NotificationEvent, 0)
	for _, event := range publisher.events {
		if event.Event == name {
			events = append(events, event)
		}
	}

	return events
}

type fakeUserRepository struct {
	store *fakeStore
}

func (repository fakeUserRepository) Create(_ context.Context, user *domain.User) error {
	repository.store.assignID(&user.ID)
	*user = repository.store.addUser(*user)
	return nil
}

func (repository fakeUserRepository) Update(_ context.Context, user *domain.User) error {
	if _, ok := repository.store.users[user.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.users[user.ID] = *user
	return nil
}

func (repository fakeUserRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.users[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.users, id)
	return nil
}

func (repository fakeUserRepository) GetByID(_ context.Context, id domain.ID) (*domain.User, error) {
	user, ok := repository.store.users[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	return &user, nil
}

func (repository fakeUserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, user := range repository.store.users {
		if user.Email == email {
			found := user
			return &found, nil
		}
	}
	return nil, application.ErrNotFound
}

func (repository fakeUserRepository) List(_ context.Context, filter application.UserFilter) ([]domain.User, error) {
	users := make([]domain.User, 0, len(repository.store.users))
	for _, user := range repository.store.users {
		if filter.Role != nil && user.Role != *filter.Role {
			continue
		}
		if filter.IsActive != nil && user.IsActive != *filter.IsActive {
			continue
		}
		if filter.Search != "" && !containsAny(user.Name, user.Email, filter.Search) {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

type fakeWarehouseRepository struct {
	store *fakeStore
}

func (repository fakeWarehouseRepository) Create(_ context.Context, warehouse *domain.Warehouse) error {
	repository.store.assignID(&warehouse.ID)
	*warehouse = repository.store.addWarehouse(*warehouse)
	return nil
}

func (repository fakeWarehouseRepository) Update(_ context.Context, warehouse *domain.Warehouse) error {
	if _, ok := repository.store.warehouses[warehouse.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.warehouses[warehouse.ID] = *warehouse
	return nil
}

func (repository fakeWarehouseRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.warehouses[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.warehouses, id)
	return nil
}

func (repository fakeWarehouseRepository) GetByID(_ context.Context, id domain.ID) (*domain.Warehouse, error) {
	warehouse, ok := repository.store.warehouses[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	return &warehouse, nil
}

func (repository fakeWarehouseRepository) GetByCode(_ context.Context, code string) (*domain.Warehouse, error) {
	for _, warehouse := range repository.store.warehouses {
		if warehouse.Code == code {
			found := warehouse
			return &found, nil
		}
	}
	return nil, application.ErrNotFound
}

func (repository fakeWarehouseRepository) List(_ context.Context, filter application.WarehouseFilter) ([]domain.Warehouse, error) {
	warehouses := make([]domain.Warehouse, 0, len(repository.store.warehouses))
	for _, warehouse := range repository.store.warehouses {
		if filter.IsActive != nil && warehouse.IsActive != *filter.IsActive {
			continue
		}
		if filter.Code != "" && warehouse.Code != filter.Code {
			continue
		}
		if filter.Search != "" && !containsAny(warehouse.Name, warehouse.Location, filter.Search) {
			continue
		}
		warehouses = append(warehouses, warehouse)
	}
	return warehouses, nil
}

func (repository fakeWarehouseRepository) SetActive(_ context.Context, id domain.ID, isActive bool) error {
	warehouse, ok := repository.store.warehouses[id]
	if !ok {
		return application.ErrNotFound
	}
	warehouse.IsActive = isActive
	repository.store.warehouses[id] = warehouse
	return nil
}

func (repository fakeWarehouseRepository) HasStock(_ context.Context, id domain.ID) (bool, error) {
	if _, ok := repository.store.warehouses[id]; !ok {
		return false, application.ErrNotFound
	}
	for _, stockItem := range repository.store.stock {
		if stockItem.WarehouseID == id && stockItem.Quantity > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (repository fakeWarehouseRepository) HasActiveOperations(_ context.Context, id domain.ID) (bool, error) {
	for _, supply := range repository.store.supplies {
		if supply.WarehouseID == id && !supply.Status.IsFinal() {
			return true, nil
		}
	}
	for _, request := range repository.store.outbounds {
		if request.WarehouseID == id && !request.Status.IsFinal() {
			return true, nil
		}
	}
	for _, transfer := range repository.store.transfers {
		if (transfer.SourceWarehouseID == id || transfer.TargetWarehouseID == id) && !transfer.Status.IsFinal() {
			return true, nil
		}
	}
	return false, nil
}

type fakeCategoryRepository struct {
	store *fakeStore
}

func (repository fakeCategoryRepository) Create(_ context.Context, category *domain.Category) error {
	repository.store.assignID(&category.ID)
	*category = repository.store.addCategory(*category)
	return nil
}

func (repository fakeCategoryRepository) Update(_ context.Context, category *domain.Category) error {
	if _, ok := repository.store.categories[category.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.categories[category.ID] = *category
	return nil
}

func (repository fakeCategoryRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.categories[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.categories, id)
	return nil
}

func (repository fakeCategoryRepository) GetByID(_ context.Context, id domain.ID) (*domain.Category, error) {
	category, ok := repository.store.categories[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	return &category, nil
}

func (repository fakeCategoryRepository) GetByName(_ context.Context, name string) (*domain.Category, error) {
	for _, category := range repository.store.categories {
		if category.Name == name {
			found := category
			return &found, nil
		}
	}
	return nil, application.ErrNotFound
}

func (repository fakeCategoryRepository) List(_ context.Context, filter application.CategoryFilter) ([]domain.Category, error) {
	categories := make([]domain.Category, 0, len(repository.store.categories))
	for _, category := range repository.store.categories {
		if filter.IsActive != nil && category.IsActive != *filter.IsActive {
			continue
		}
		if filter.Search != "" && !containsAny(category.Name, category.Description, filter.Search) {
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (repository fakeCategoryRepository) HasActiveProducts(_ context.Context, id domain.ID) (bool, error) {
	for _, product := range repository.store.products {
		if product.CategoryID == id && product.IsActive {
			return true, nil
		}
	}
	return false, nil
}

type fakeProductRepository struct {
	store *fakeStore
}

func (repository fakeProductRepository) Create(_ context.Context, product *domain.Product) error {
	repository.store.assignID(&product.ID)
	*product = repository.store.addProduct(*product)
	return nil
}

func (repository fakeProductRepository) Update(_ context.Context, product *domain.Product) error {
	if _, ok := repository.store.products[product.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.products[product.ID] = *product
	return nil
}

func (repository fakeProductRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.products[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.products, id)
	return nil
}

func (repository fakeProductRepository) GetByID(_ context.Context, id domain.ID) (*domain.Product, error) {
	product, ok := repository.store.products[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	return &product, nil
}

func (repository fakeProductRepository) GetBySKU(_ context.Context, sku string) (*domain.Product, error) {
	for _, product := range repository.store.products {
		if product.SKU == sku {
			found := product
			return &found, nil
		}
	}
	return nil, application.ErrNotFound
}

func (repository fakeProductRepository) List(_ context.Context, filter application.ProductFilter) ([]domain.Product, error) {
	products := make([]domain.Product, 0, len(repository.store.products))
	for _, product := range repository.store.products {
		if filter.CategoryID != nil && product.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.IsActive != nil && product.IsActive != *filter.IsActive {
			continue
		}
		if filter.SKU != "" && product.SKU != filter.SKU {
			continue
		}
		if filter.Name != "" && product.Name != filter.Name {
			continue
		}
		if filter.Search != "" && !containsAny(product.Name, product.SKU, filter.Search) {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func (repository fakeProductRepository) SetActive(_ context.Context, id domain.ID, isActive bool) error {
	product, ok := repository.store.products[id]
	if !ok {
		return application.ErrNotFound
	}
	product.IsActive = isActive
	repository.store.products[id] = product
	return nil
}

func (repository fakeProductRepository) HasConfirmedMovements(_ context.Context, id domain.ID) (bool, error) {
	for _, supply := range repository.store.supplies {
		if supply.Status != domain.SupplyStatusCompleted {
			continue
		}
		for _, item := range repository.store.supplyItemsFor(supply.ID) {
			if item.ProductID == id {
				return true, nil
			}
		}
	}
	for _, request := range repository.store.outbounds {
		if request.Status != domain.OutboundRequestStatusCompleted {
			continue
		}
		for _, item := range repository.store.outboundItemsFor(request.ID) {
			if item.ProductID == id {
				return true, nil
			}
		}
	}
	for _, transfer := range repository.store.transfers {
		if transfer.Status != domain.TransferStatusCompleted {
			continue
		}
		for _, item := range repository.store.transferItemsFor(transfer.ID) {
			if item.ProductID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

type fakeStockRepository struct {
	store *fakeStore
}

func (repository fakeStockRepository) GetByWarehouse(_ context.Context, warehouseID domain.ID) ([]domain.StockItem, error) {
	stockItems := make([]domain.StockItem, 0)
	for _, stockItem := range repository.store.stock {
		if stockItem.WarehouseID == warehouseID {
			stockItems = append(stockItems, stockItem)
		}
	}
	return stockItems, nil
}

func (repository fakeStockRepository) GetByProductAcrossWarehouses(_ context.Context, productID domain.ID) ([]domain.StockItem, error) {
	stockItems := make([]domain.StockItem, 0)
	for _, stockItem := range repository.store.stock {
		if stockItem.ProductID == productID {
			stockItems = append(stockItems, stockItem)
		}
	}
	return stockItems, nil
}

func (repository fakeStockRepository) GetByWarehouseAndProduct(_ context.Context, warehouseID domain.ID, productID domain.ID) (*domain.StockItem, error) {
	stockItem, ok := repository.store.stock[stockKey{warehouseID: warehouseID, productID: productID}]
	if !ok {
		return nil, application.ErrNotFound
	}
	return &stockItem, nil
}

func (repository fakeStockRepository) Upsert(_ context.Context, stockItem *domain.StockItem) error {
	key := stockKey{warehouseID: stockItem.WarehouseID, productID: stockItem.ProductID}
	existing, ok := repository.store.stock[key]
	if ok {
		existing.Quantity = stockItem.Quantity
		repository.store.stock[key] = existing
		*stockItem = existing
		return nil
	}

	repository.store.assignID(&stockItem.ID)
	repository.store.stock[key] = *stockItem
	return nil
}

func (repository fakeStockRepository) Increment(_ context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error) {
	key := stockKey{warehouseID: warehouseID, productID: productID}
	stockItem, ok := repository.store.stock[key]
	if !ok {
		stockItem = domain.StockItem{WarehouseID: warehouseID, ProductID: productID}
		repository.store.assignID(&stockItem.ID)
	}
	stockItem.Quantity += quantity
	repository.store.stock[key] = stockItem
	return &stockItem, nil
}

func (repository fakeStockRepository) Decrement(_ context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error) {
	key := stockKey{warehouseID: warehouseID, productID: productID}
	stockItem, ok := repository.store.stock[key]
	if !ok || stockItem.Quantity < quantity {
		return nil, application.ErrNotFound
	}
	stockItem.Quantity -= quantity
	repository.store.stock[key] = stockItem
	return &stockItem, nil
}

func (repository fakeStockRepository) GetTotalQuantityForProduct(_ context.Context, productID domain.ID) (int64, error) {
	var total int64
	for _, stockItem := range repository.store.stock {
		if stockItem.ProductID == productID {
			total += stockItem.Quantity
		}
	}
	return total, nil
}

func (repository fakeStockRepository) GetLowStock(_ context.Context) ([]domain.StockItem, error) {
	stockItems := make([]domain.StockItem, 0)
	for _, stockItem := range repository.store.stock {
		product, ok := repository.store.products[stockItem.ProductID]
		if !ok || !product.IsActive {
			continue
		}
		if stockItem.Quantity <= product.MinStockThreshold {
			stockItems = append(stockItems, stockItem)
		}
	}
	return stockItems, nil
}

type fakeSupplyRepository struct {
	store *fakeStore
}

func (repository fakeSupplyRepository) Create(_ context.Context, supply *domain.Supply) error {
	repository.store.assignID(&supply.ID)
	repository.store.supplies[supply.ID] = *supply
	return nil
}

func (repository fakeSupplyRepository) Update(_ context.Context, supply *domain.Supply) error {
	if _, ok := repository.store.supplies[supply.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.supplies[supply.ID] = *supply
	return nil
}

func (repository fakeSupplyRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.supplies[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.supplies, id)
	for itemID, item := range repository.store.supplyItems {
		if item.SupplyID == id {
			delete(repository.store.supplyItems, itemID)
		}
	}
	return nil
}

func (repository fakeSupplyRepository) GetByID(_ context.Context, id domain.ID) (*domain.Supply, error) {
	supply, ok := repository.store.supplies[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	supply.Items = repository.store.supplyItemsFor(id)
	return &supply, nil
}

func (repository fakeSupplyRepository) List(_ context.Context, filter application.SupplyFilter) ([]domain.Supply, error) {
	supplies := make([]domain.Supply, 0, len(repository.store.supplies))
	for _, supply := range repository.store.supplies {
		if filter.WarehouseID != nil && supply.WarehouseID != *filter.WarehouseID {
			continue
		}
		if filter.Status != nil && supply.Status != *filter.Status {
			continue
		}
		supply.Items = repository.store.supplyItemsFor(supply.ID)
		supplies = append(supplies, supply)
	}
	return supplies, nil
}

func (repository fakeSupplyRepository) AddItem(_ context.Context, item *domain.SupplyItem) error {
	repository.store.assignID(&item.ID)
	repository.store.supplyItems[item.ID] = *item
	return nil
}

func (repository fakeSupplyRepository) UpdateItem(_ context.Context, item *domain.SupplyItem) error {
	if _, ok := repository.store.supplyItems[item.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.supplyItems[item.ID] = *item
	return nil
}

func (repository fakeSupplyRepository) DeleteItem(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.supplyItems[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.supplyItems, id)
	return nil
}

func (repository fakeSupplyRepository) ListItems(_ context.Context, supplyID domain.ID) ([]domain.SupplyItem, error) {
	return repository.store.supplyItemsFor(supplyID), nil
}

type fakeOutboundRepository struct {
	store *fakeStore
}

func (repository fakeOutboundRepository) Create(_ context.Context, request *domain.OutboundRequest) error {
	repository.store.assignID(&request.ID)
	repository.store.outbounds[request.ID] = *request
	return nil
}

func (repository fakeOutboundRepository) Update(_ context.Context, request *domain.OutboundRequest) error {
	if _, ok := repository.store.outbounds[request.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.outbounds[request.ID] = *request
	return nil
}

func (repository fakeOutboundRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.outbounds[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.outbounds, id)
	for itemID, item := range repository.store.outboundItems {
		if item.OutboundRequestID == id {
			delete(repository.store.outboundItems, itemID)
		}
	}
	return nil
}

func (repository fakeOutboundRepository) GetByID(_ context.Context, id domain.ID) (*domain.OutboundRequest, error) {
	request, ok := repository.store.outbounds[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	request.Items = repository.store.outboundItemsFor(id)
	return &request, nil
}

func (repository fakeOutboundRepository) List(_ context.Context, filter application.OutboundRequestFilter) ([]domain.OutboundRequest, error) {
	requests := make([]domain.OutboundRequest, 0, len(repository.store.outbounds))
	for _, request := range repository.store.outbounds {
		if filter.WarehouseID != nil && request.WarehouseID != *filter.WarehouseID {
			continue
		}
		if filter.Status != nil && request.Status != *filter.Status {
			continue
		}
		request.Items = repository.store.outboundItemsFor(request.ID)
		requests = append(requests, request)
	}
	return requests, nil
}

func (repository fakeOutboundRepository) AddItem(_ context.Context, item *domain.OutboundRequestItem) error {
	repository.store.assignID(&item.ID)
	repository.store.outboundItems[item.ID] = *item
	return nil
}

func (repository fakeOutboundRepository) UpdateItem(_ context.Context, item *domain.OutboundRequestItem) error {
	if _, ok := repository.store.outboundItems[item.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.outboundItems[item.ID] = *item
	return nil
}

func (repository fakeOutboundRepository) DeleteItem(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.outboundItems[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.outboundItems, id)
	return nil
}

func (repository fakeOutboundRepository) ListItems(_ context.Context, outboundRequestID domain.ID) ([]domain.OutboundRequestItem, error) {
	return repository.store.outboundItemsFor(outboundRequestID), nil
}

type fakeTransferRepository struct {
	store *fakeStore
}

func (repository fakeTransferRepository) Create(_ context.Context, transfer *domain.Transfer) error {
	repository.store.assignID(&transfer.ID)
	repository.store.transfers[transfer.ID] = *transfer
	return nil
}

func (repository fakeTransferRepository) Update(_ context.Context, transfer *domain.Transfer) error {
	if _, ok := repository.store.transfers[transfer.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.transfers[transfer.ID] = *transfer
	return nil
}

func (repository fakeTransferRepository) Delete(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.transfers[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.transfers, id)
	for itemID, item := range repository.store.transferItems {
		if item.TransferID == id {
			delete(repository.store.transferItems, itemID)
		}
	}
	return nil
}

func (repository fakeTransferRepository) GetByID(_ context.Context, id domain.ID) (*domain.Transfer, error) {
	transfer, ok := repository.store.transfers[id]
	if !ok {
		return nil, application.ErrNotFound
	}
	transfer.Items = repository.store.transferItemsFor(id)
	return &transfer, nil
}

func (repository fakeTransferRepository) List(_ context.Context, filter application.TransferFilter) ([]domain.Transfer, error) {
	transfers := make([]domain.Transfer, 0, len(repository.store.transfers))
	for _, transfer := range repository.store.transfers {
		if filter.WarehouseID != nil && transfer.SourceWarehouseID != *filter.WarehouseID && transfer.TargetWarehouseID != *filter.WarehouseID {
			continue
		}
		if filter.SourceWarehouseID != nil && transfer.SourceWarehouseID != *filter.SourceWarehouseID {
			continue
		}
		if filter.TargetWarehouseID != nil && transfer.TargetWarehouseID != *filter.TargetWarehouseID {
			continue
		}
		if filter.Status != nil && transfer.Status != *filter.Status {
			continue
		}
		transfer.Items = repository.store.transferItemsFor(transfer.ID)
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func (repository fakeTransferRepository) AddItem(_ context.Context, item *domain.TransferItem) error {
	repository.store.assignID(&item.ID)
	repository.store.transferItems[item.ID] = *item
	return nil
}

func (repository fakeTransferRepository) UpdateItem(_ context.Context, item *domain.TransferItem) error {
	if _, ok := repository.store.transferItems[item.ID]; !ok {
		return application.ErrNotFound
	}
	repository.store.transferItems[item.ID] = *item
	return nil
}

func (repository fakeTransferRepository) DeleteItem(_ context.Context, id domain.ID) error {
	if _, ok := repository.store.transferItems[id]; !ok {
		return application.ErrNotFound
	}
	delete(repository.store.transferItems, id)
	return nil
}

func (repository fakeTransferRepository) ListItems(_ context.Context, transferID domain.ID) ([]domain.TransferItem, error) {
	return repository.store.transferItemsFor(transferID), nil
}

type fakeMovementRepository struct {
	store *fakeStore
}

func (repository fakeMovementRepository) ListByProduct(_ context.Context, productID domain.ID, _ application.ListOptions) ([]application.MovementRecord, error) {
	records := make([]application.MovementRecord, 0)
	for _, record := range repository.store.movementRecords() {
		if record.ProductID == productID {
			records = append(records, record)
		}
	}
	return records, nil
}

func (repository fakeMovementRepository) ListByWarehouse(_ context.Context, warehouseID domain.ID, _ application.ListOptions) ([]application.MovementRecord, error) {
	records := make([]application.MovementRecord, 0)
	for _, record := range repository.store.movementRecords() {
		if record.WarehouseID == warehouseID {
			records = append(records, record)
		}
	}
	return records, nil
}

func (store *fakeStore) supplyItemsFor(supplyID domain.ID) []domain.SupplyItem {
	items := make([]domain.SupplyItem, 0)
	for _, item := range store.supplyItems {
		if item.SupplyID == supplyID {
			items = append(items, item)
		}
	}
	return items
}

func (store *fakeStore) outboundItemsFor(requestID domain.ID) []domain.OutboundRequestItem {
	items := make([]domain.OutboundRequestItem, 0)
	for _, item := range store.outboundItems {
		if item.OutboundRequestID == requestID {
			items = append(items, item)
		}
	}
	return items
}

func (store *fakeStore) transferItemsFor(transferID domain.ID) []domain.TransferItem {
	items := make([]domain.TransferItem, 0)
	for _, item := range store.transferItems {
		if item.TransferID == transferID {
			items = append(items, item)
		}
	}
	return items
}

func (store *fakeStore) movementRecords() []application.MovementRecord {
	records := make([]application.MovementRecord, 0)
	for _, supply := range store.supplies {
		if supply.Status != domain.SupplyStatusCompleted {
			continue
		}
		for _, item := range store.supplyItemsFor(supply.ID) {
			records = append(records, application.MovementRecord{
				Kind:            application.MovementKindSupply,
				OperationID:     supply.ID,
				OperationItemID: item.ID,
				ProductID:       item.ProductID,
				WarehouseID:     supply.WarehouseID,
				Quantity:        item.Quantity,
				Status:          string(supply.Status),
				OccurredAt:      defaultTime(time.Time{}),
			})
		}
	}
	for _, request := range store.outbounds {
		if request.Status != domain.OutboundRequestStatusCompleted {
			continue
		}
		for _, item := range store.outboundItemsFor(request.ID) {
			records = append(records, application.MovementRecord{
				Kind:            application.MovementKindOutbound,
				OperationID:     request.ID,
				OperationItemID: item.ID,
				ProductID:       item.ProductID,
				WarehouseID:     request.WarehouseID,
				Quantity:        -item.Quantity,
				Status:          string(request.Status),
				OccurredAt:      defaultTime(time.Time{}),
			})
		}
	}
	for _, transfer := range store.transfers {
		if transfer.Status != domain.TransferStatusCompleted {
			continue
		}
		for _, item := range store.transferItemsFor(transfer.ID) {
			sourceWarehouseID := transfer.TargetWarehouseID
			targetWarehouseID := transfer.SourceWarehouseID
			records = append(records, application.MovementRecord{
				Kind:               application.MovementKindTransferOut,
				OperationID:        transfer.ID,
				OperationItemID:    item.ID,
				ProductID:          item.ProductID,
				WarehouseID:        transfer.SourceWarehouseID,
				RelatedWarehouseID: &sourceWarehouseID,
				Quantity:           -item.Quantity,
				Status:             string(transfer.Status),
				OccurredAt:         defaultTime(time.Time{}),
			})
			records = append(records, application.MovementRecord{
				Kind:               application.MovementKindTransferIn,
				OperationID:        transfer.ID,
				OperationItemID:    item.ID,
				ProductID:          item.ProductID,
				WarehouseID:        transfer.TargetWarehouseID,
				RelatedWarehouseID: &targetWarehouseID,
				Quantity:           item.Quantity,
				Status:             string(transfer.Status),
				OccurredAt:         defaultTime(time.Time{}),
			})
		}
	}
	return records
}

func containsAny(left string, right string, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(left), query) || strings.Contains(strings.ToLower(right), query)
}
