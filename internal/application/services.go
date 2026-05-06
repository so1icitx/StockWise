package application

// Services groups all application services.
type Services struct {
	Users            *UserService
	Warehouses       *WarehouseService
	Categories       *CategoryService
	Products         *ProductService
	Stock            *StockService
	Supplies         *SupplyService
	OutboundRequests *OutboundRequestService
	Transfers        *TransferService
	Movements        *MovementService
	Notifications    NotificationPublisher
}

// NewServices creates all application services from repository and transaction dependencies.
func NewServices(provider RepositoryProvider, transactions TransactionManager, publishers ...NotificationPublisher) Services {
	notifications := notificationPublisherFrom(publishers...)

	return Services{
		Users:            NewUserService(provider),
		Warehouses:       NewWarehouseService(provider, notifications),
		Categories:       NewCategoryService(provider),
		Products:         NewProductService(provider),
		Stock:            NewStockService(provider),
		Supplies:         NewSupplyService(provider, transactions, notifications),
		OutboundRequests: NewOutboundRequestService(provider, transactions, notifications),
		Transfers:        NewTransferService(provider, transactions, notifications),
		Movements:        NewMovementService(provider),
		Notifications:    notifications,
	}
}
