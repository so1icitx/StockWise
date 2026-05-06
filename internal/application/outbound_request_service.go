package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// OutboundRequestService handles outbound request workflows.
type OutboundRequestService struct {
	provider      RepositoryProvider
	transactions  TransactionManager
	notifications NotificationPublisher
}

// NewOutboundRequestService creates an outbound request service.
func NewOutboundRequestService(provider RepositoryProvider, transactions TransactionManager, publishers ...NotificationPublisher) *OutboundRequestService {
	return &OutboundRequestService{
		provider:      provider,
		transactions:  transactions,
		notifications: notificationPublisherFrom(publishers...),
	}
}

// Create creates a draft outbound request for an active warehouse.
func (service *OutboundRequestService) Create(ctx context.Context, input CreateOutboundRequestInput) (*domain.OutboundRequest, error) {
	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, input.CreatedByUserID); err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, input.WarehouseID); err != nil {
		return nil, err
	}

	request := domain.OutboundRequest{
		WarehouseID:     input.WarehouseID,
		Status:          domain.OutboundRequestStatusDraft,
		CreatedByUserID: input.CreatedByUserID,
	}
	if err := repos.OutboundRequests.Create(ctx, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

// AddItem adds a product row to an editable outbound request.
func (service *OutboundRequestService) AddItem(ctx context.Context, requestID domain.ID, input OutboundRequestItemInput) (*domain.OutboundRequestItem, error) {
	if err := requireID(requestID, "outbound request id"); err != nil {
		return nil, err
	}
	if err := validateOutboundItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	request, err := service.requireEditableRequest(ctx, repos, requestID)
	if err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, request.WarehouseID); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureOutboundProductNotDuplicated(ctx, repos.OutboundRequests, requestID, 0, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.OutboundRequestItem{
		OutboundRequestID: requestID,
		ProductID:         input.ProductID,
		Quantity:          input.Quantity,
	}
	if err := repos.OutboundRequests.AddItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItem updates a product row on an editable outbound request.
func (service *OutboundRequestService) UpdateItem(ctx context.Context, requestID domain.ID, itemID domain.ID, input OutboundRequestItemInput) (*domain.OutboundRequestItem, error) {
	if err := requireID(requestID, "outbound request id"); err != nil {
		return nil, err
	}
	if err := requireID(itemID, "outbound request item id"); err != nil {
		return nil, err
	}
	if err := validateOutboundItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	request, err := service.requireEditableRequest(ctx, repos, requestID)
	if err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, request.WarehouseID); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureOutboundProductNotDuplicated(ctx, repos.OutboundRequests, requestID, itemID, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.OutboundRequestItem{
		ID:                itemID,
		OutboundRequestID: requestID,
		ProductID:         input.ProductID,
		Quantity:          input.Quantity,
	}
	if err := repos.OutboundRequests.UpdateItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem removes a product row from an editable outbound request.
func (service *OutboundRequestService) DeleteItem(ctx context.Context, requestID domain.ID, itemID domain.ID) error {
	if err := requireID(requestID, "outbound request id"); err != nil {
		return err
	}
	if err := requireID(itemID, "outbound request item id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableRequest(ctx, repos, requestID); err != nil {
		return err
	}
	if err := ensureOutboundItemBelongs(ctx, repos.OutboundRequests, requestID, itemID); err != nil {
		return err
	}

	return repos.OutboundRequests.DeleteItem(ctx, itemID)
}

// Approve approves a draft outbound request with at least one item.
func (service *OutboundRequestService) Approve(ctx context.Context, id domain.ID, approvedByUserID domain.ID) (*domain.OutboundRequest, error) {
	if err := requireID(id, "outbound request id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, approvedByUserID); err != nil {
		return nil, err
	}
	request, err := repos.OutboundRequests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !request.Status.CanApprove() {
		return nil, businessRuleError(ErrOperationLocked, "outbound request cannot be approved in its current status")
	}
	if !request.HasItems() {
		return nil, businessRuleError(ErrBusinessRule, "cannot approve outbound request without items")
	}
	if !request.HasValidItems() {
		return nil, validationError("outbound request items must have positive quantities")
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, request.WarehouseID); err != nil {
		return nil, err
	}
	for _, item := range request.Items {
		if _, err := requireActiveProduct(ctx, repos.Products, item.ProductID); err != nil {
			return nil, err
		}
	}

	now := nowUTC()
	request.Status = domain.OutboundRequestStatusApproved
	request.ApprovedByUserID = &approvedByUserID
	request.ApprovedAt = &now
	if err := repos.OutboundRequests.Update(ctx, request); err != nil {
		return nil, err
	}

	approved, err := repos.OutboundRequests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	service.notifications.Publish(ctx, outboundApprovedEvent(*approved))

	return approved, nil
}

// Execute executes an approved outbound request and decreases stock inside one transaction.
func (service *OutboundRequestService) Execute(ctx context.Context, id domain.ID, executedByUserID domain.ID) (*domain.OutboundRequest, error) {
	if err := requireID(id, "outbound request id"); err != nil {
		return nil, err
	}
	if service.transactions == nil {
		return nil, businessRuleError(ErrBusinessRule, "transaction manager is required")
	}

	var completed *domain.OutboundRequest
	var stockEvents []NotificationEvent
	err := service.transactions.WithinTransaction(ctx, func(txCtx context.Context, repos Repositories) error {
		if _, err := requireActiveUser(txCtx, repos.Users, executedByUserID); err != nil {
			return err
		}
		request, err := repos.OutboundRequests.GetByID(txCtx, id)
		if err != nil {
			return err
		}
		if !request.Status.CanExecute() {
			return businessRuleError(ErrOperationLocked, "only approved outbound requests can be executed")
		}
		if !request.HasItems() {
			return businessRuleError(ErrBusinessRule, "cannot execute outbound request without items")
		}
		if !request.HasValidItems() {
			return validationError("outbound request items must have positive quantities")
		}
		if _, err := requireActiveWarehouse(txCtx, repos.Warehouses, request.WarehouseID); err != nil {
			return err
		}

		products := make(map[domain.ID]domain.Product, len(request.Items))
		for _, item := range request.Items {
			product, err := requireActiveProduct(txCtx, repos.Products, item.ProductID)
			if err != nil {
				return err
			}
			products[item.ProductID] = *product
			stockItem, err := repos.Stock.GetByWarehouseAndProduct(txCtx, request.WarehouseID, item.ProductID)
			if err != nil {
				if isNotFound(err) {
					return businessRuleError(ErrInsufficientStock, "not enough stock to execute outbound request")
				}
				return err
			}
			if stockItem.Quantity < item.Quantity {
				return businessRuleError(ErrInsufficientStock, "not enough stock to execute outbound request")
			}
		}

		for _, item := range request.Items {
			stockItem, err := repos.Stock.Decrement(txCtx, request.WarehouseID, item.ProductID, item.Quantity)
			if err != nil {
				if isNotFound(err) {
					return businessRuleError(ErrInsufficientStock, "not enough stock to execute outbound request")
				}
				return err
			}
			stockEvents = appendStockStateEvent(stockEvents, *stockItem, products[item.ProductID], NotificationOutboundCompleted)
		}

		now := nowUTC()
		request.Status = domain.OutboundRequestStatusCompleted
		request.ExecutedByUserID = &executedByUserID
		request.ExecutedAt = &now
		if err := repos.OutboundRequests.Update(txCtx, request); err != nil {
			return err
		}

		completed, err = repos.OutboundRequests.GetByID(txCtx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	service.notifications.Publish(ctx, outboundCompletedEvent(*completed))
	publishAll(ctx, service.notifications, stockEvents)

	return completed, nil
}

// Cancel cancels a draft or approved outbound request.
func (service *OutboundRequestService) Cancel(ctx context.Context, id domain.ID, cancelledByUserID domain.ID) (*domain.OutboundRequest, error) {
	if err := requireID(id, "outbound request id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, cancelledByUserID); err != nil {
		return nil, err
	}
	request, err := repos.OutboundRequests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !request.Status.CanCancel() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled outbound request cannot be cancelled")
	}

	now := nowUTC()
	request.Status = domain.OutboundRequestStatusCancelled
	request.CancelledByUserID = &cancelledByUserID
	request.CancelledAt = &now
	if err := repos.OutboundRequests.Update(ctx, request); err != nil {
		return nil, err
	}

	return repos.OutboundRequests.GetByID(ctx, id)
}

// Reject cancels an outbound request as a rejected request.
func (service *OutboundRequestService) Reject(ctx context.Context, id domain.ID, cancelledByUserID domain.ID) (*domain.OutboundRequest, error) {
	return service.Cancel(ctx, id, cancelledByUserID)
}

// Delete removes a draft outbound request.
func (service *OutboundRequestService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "outbound request id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableRequest(ctx, repos, id); err != nil {
		return err
	}

	return repos.OutboundRequests.Delete(ctx, id)
}

// GetByID returns an outbound request by identifier.
func (service *OutboundRequestService) GetByID(ctx context.Context, id domain.ID) (*domain.OutboundRequest, error) {
	if err := requireID(id, "outbound request id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).OutboundRequests.GetByID(ctx, id)
}

// List returns outbound requests matching the filter.
func (service *OutboundRequestService) List(ctx context.Context, filter OutboundRequestFilter) ([]domain.OutboundRequest, error) {
	return repositories(service.provider).OutboundRequests.List(ctx, filter)
}

func (service *OutboundRequestService) requireEditableRequest(ctx context.Context, repos Repositories, id domain.ID) (*domain.OutboundRequest, error) {
	request, err := repos.OutboundRequests.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !request.CanBeEdited() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled outbound request cannot be edited")
	}

	return request, nil
}

func validateOutboundItemInput(input OutboundRequestItemInput) error {
	if err := requireID(input.ProductID, "product id"); err != nil {
		return err
	}
	return ensurePositiveQuantity(input.Quantity)
}

func ensureOutboundProductNotDuplicated(ctx context.Context, repository OutboundRequestRepository, requestID domain.ID, currentItemID domain.ID, productID domain.ID) error {
	items, err := repository.ListItems(ctx, requestID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ProductID == productID && item.ID != currentItemID {
			return conflictError(ErrConflict, "product already exists in outbound request")
		}
	}

	return nil
}

func ensureOutboundItemBelongs(ctx context.Context, repository OutboundRequestRepository, requestID domain.ID, itemID domain.ID) error {
	items, err := repository.ListItems(ctx, requestID)
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
