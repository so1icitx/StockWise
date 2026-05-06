package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// TransferService handles stock transfer workflows.
type TransferService struct {
	provider      RepositoryProvider
	transactions  TransactionManager
	notifications NotificationPublisher
}

// NewTransferService creates a transfer service.
func NewTransferService(provider RepositoryProvider, transactions TransactionManager, publishers ...NotificationPublisher) *TransferService {
	return &TransferService{
		provider:      provider,
		transactions:  transactions,
		notifications: notificationPublisherFrom(publishers...),
	}
}

// Create creates a draft transfer between two active, different warehouses.
func (service *TransferService) Create(ctx context.Context, input CreateTransferInput) (*domain.Transfer, error) {
	if input.SourceWarehouseID == input.TargetWarehouseID && !input.SourceWarehouseID.IsZero() {
		return nil, businessRuleError(ErrBusinessRule, "cannot transfer to the same warehouse")
	}

	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, input.CreatedByUserID); err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, input.SourceWarehouseID); err != nil {
		return nil, err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, input.TargetWarehouseID); err != nil {
		return nil, err
	}

	transfer := domain.Transfer{
		SourceWarehouseID: input.SourceWarehouseID,
		TargetWarehouseID: input.TargetWarehouseID,
		Status:            domain.TransferStatusDraft,
		CreatedByUserID:   input.CreatedByUserID,
	}
	if err := repos.Transfers.Create(ctx, &transfer); err != nil {
		return nil, err
	}

	return &transfer, nil
}

// AddItem adds a product row to an editable transfer.
func (service *TransferService) AddItem(ctx context.Context, transferID domain.ID, input TransferItemInput) (*domain.TransferItem, error) {
	if err := requireID(transferID, "transfer id"); err != nil {
		return nil, err
	}
	if err := validateTransferItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	transfer, err := service.requireEditableTransfer(ctx, repos, transferID)
	if err != nil {
		return nil, err
	}
	if err := ensureTransferWarehousesUsable(ctx, repos, transfer); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureTransferProductNotDuplicated(ctx, repos.Transfers, transferID, 0, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.TransferItem{
		TransferID: transferID,
		ProductID:  input.ProductID,
		Quantity:   input.Quantity,
	}
	if err := repos.Transfers.AddItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItem updates a product row on an editable transfer.
func (service *TransferService) UpdateItem(ctx context.Context, transferID domain.ID, itemID domain.ID, input TransferItemInput) (*domain.TransferItem, error) {
	if err := requireID(transferID, "transfer id"); err != nil {
		return nil, err
	}
	if err := requireID(itemID, "transfer item id"); err != nil {
		return nil, err
	}
	if err := validateTransferItemInput(input); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	transfer, err := service.requireEditableTransfer(ctx, repos, transferID)
	if err != nil {
		return nil, err
	}
	if err := ensureTransferWarehousesUsable(ctx, repos, transfer); err != nil {
		return nil, err
	}
	if _, err := requireActiveProduct(ctx, repos.Products, input.ProductID); err != nil {
		return nil, err
	}
	if err := ensureTransferProductNotDuplicated(ctx, repos.Transfers, transferID, itemID, input.ProductID); err != nil {
		return nil, err
	}

	item := domain.TransferItem{
		ID:         itemID,
		TransferID: transferID,
		ProductID:  input.ProductID,
		Quantity:   input.Quantity,
	}
	if err := repos.Transfers.UpdateItem(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem removes a product row from an editable transfer.
func (service *TransferService) DeleteItem(ctx context.Context, transferID domain.ID, itemID domain.ID) error {
	if err := requireID(transferID, "transfer id"); err != nil {
		return err
	}
	if err := requireID(itemID, "transfer item id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableTransfer(ctx, repos, transferID); err != nil {
		return err
	}
	if err := ensureTransferItemBelongs(ctx, repos.Transfers, transferID, itemID); err != nil {
		return err
	}

	return repos.Transfers.DeleteItem(ctx, itemID)
}

// Confirm confirms a transfer and moves stock inside one transaction.
func (service *TransferService) Confirm(ctx context.Context, id domain.ID, confirmedByUserID domain.ID) (*domain.Transfer, error) {
	if err := requireID(id, "transfer id"); err != nil {
		return nil, err
	}
	if service.transactions == nil {
		return nil, businessRuleError(ErrBusinessRule, "transaction manager is required")
	}

	var completed *domain.Transfer
	var stockEvents []NotificationEvent
	err := service.transactions.WithinTransaction(ctx, func(txCtx context.Context, repos Repositories) error {
		if _, err := requireActiveUser(txCtx, repos.Users, confirmedByUserID); err != nil {
			return err
		}
		transfer, err := repos.Transfers.GetByID(txCtx, id)
		if err != nil {
			return err
		}
		if !transfer.Status.CanConfirm() {
			return businessRuleError(ErrOperationLocked, "completed or cancelled transfer cannot be confirmed")
		}
		if transfer.UsesSameWarehouse() {
			return businessRuleError(ErrBusinessRule, "cannot transfer to the same warehouse")
		}
		if !transfer.HasItems() {
			return businessRuleError(ErrBusinessRule, "cannot confirm transfer without items")
		}
		if !transfer.HasValidItems() {
			return validationError("transfer items must have positive quantities")
		}
		if err := ensureTransferWarehousesUsable(txCtx, repos, transfer); err != nil {
			return err
		}

		products := make(map[domain.ID]domain.Product, len(transfer.Items))
		for _, item := range transfer.Items {
			product, err := requireActiveProduct(txCtx, repos.Products, item.ProductID)
			if err != nil {
				return err
			}
			products[item.ProductID] = *product
			stockItem, err := repos.Stock.GetByWarehouseAndProduct(txCtx, transfer.SourceWarehouseID, item.ProductID)
			if err != nil {
				if isNotFound(err) {
					return businessRuleError(ErrInsufficientStock, "not enough stock to confirm transfer")
				}
				return err
			}
			if stockItem.Quantity < item.Quantity {
				return businessRuleError(ErrInsufficientStock, "not enough stock to confirm transfer")
			}
		}

		for _, item := range transfer.Items {
			sourceStockItem, err := repos.Stock.Decrement(txCtx, transfer.SourceWarehouseID, item.ProductID, item.Quantity)
			if err != nil {
				if isNotFound(err) {
					return businessRuleError(ErrInsufficientStock, "not enough stock to confirm transfer")
				}
				return err
			}
			stockEvents = appendStockStateEvent(stockEvents, *sourceStockItem, products[item.ProductID], NotificationTransferConfirmed)

			targetStockItem, err := repos.Stock.Increment(txCtx, transfer.TargetWarehouseID, item.ProductID, item.Quantity)
			if err != nil {
				return err
			}
			stockEvents = appendStockStateEvent(stockEvents, *targetStockItem, products[item.ProductID], NotificationTransferConfirmed)
		}

		now := nowUTC()
		transfer.Status = domain.TransferStatusCompleted
		transfer.ConfirmedByUserID = &confirmedByUserID
		transfer.ConfirmedAt = &now
		if err := repos.Transfers.Update(txCtx, transfer); err != nil {
			return err
		}

		completed, err = repos.Transfers.GetByID(txCtx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	service.notifications.Publish(ctx, transferConfirmedEvent(*completed))
	publishAll(ctx, service.notifications, stockEvents)

	return completed, nil
}

// Cancel cancels a draft transfer without moving stock.
func (service *TransferService) Cancel(ctx context.Context, id domain.ID, cancelledByUserID domain.ID) (*domain.Transfer, error) {
	if err := requireID(id, "transfer id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := requireActiveUser(ctx, repos.Users, cancelledByUserID); err != nil {
		return nil, err
	}
	transfer, err := repos.Transfers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !transfer.Status.CanCancel() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled transfer cannot be cancelled")
	}

	now := nowUTC()
	transfer.Status = domain.TransferStatusCancelled
	transfer.CancelledByUserID = &cancelledByUserID
	transfer.CancelledAt = &now
	if err := repos.Transfers.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return repos.Transfers.GetByID(ctx, id)
}

// Delete removes a draft transfer.
func (service *TransferService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "transfer id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := service.requireEditableTransfer(ctx, repos, id); err != nil {
		return err
	}

	return repos.Transfers.Delete(ctx, id)
}

// GetByID returns a transfer by identifier.
func (service *TransferService) GetByID(ctx context.Context, id domain.ID) (*domain.Transfer, error) {
	if err := requireID(id, "transfer id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Transfers.GetByID(ctx, id)
}

// List returns transfers matching the filter.
func (service *TransferService) List(ctx context.Context, filter TransferFilter) ([]domain.Transfer, error) {
	return repositories(service.provider).Transfers.List(ctx, filter)
}

func (service *TransferService) requireEditableTransfer(ctx context.Context, repos Repositories, id domain.ID) (*domain.Transfer, error) {
	transfer, err := repos.Transfers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !transfer.CanBeEdited() {
		return nil, businessRuleError(ErrOperationLocked, "completed or cancelled transfer cannot be edited")
	}

	return transfer, nil
}

func validateTransferItemInput(input TransferItemInput) error {
	if err := requireID(input.ProductID, "product id"); err != nil {
		return err
	}
	return ensurePositiveQuantity(input.Quantity)
}

func ensureTransferWarehousesUsable(ctx context.Context, repos Repositories, transfer *domain.Transfer) error {
	if transfer.UsesSameWarehouse() {
		return businessRuleError(ErrBusinessRule, "cannot transfer to the same warehouse")
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, transfer.SourceWarehouseID); err != nil {
		return err
	}
	if _, err := requireActiveWarehouse(ctx, repos.Warehouses, transfer.TargetWarehouseID); err != nil {
		return err
	}

	return nil
}

func ensureTransferProductNotDuplicated(ctx context.Context, repository TransferRepository, transferID domain.ID, currentItemID domain.ID, productID domain.ID) error {
	items, err := repository.ListItems(ctx, transferID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ProductID == productID && item.ID != currentItemID {
			return conflictError(ErrConflict, "product already exists in transfer")
		}
	}

	return nil
}

func ensureTransferItemBelongs(ctx context.Context, repository TransferRepository, transferID domain.ID, itemID domain.ID) error {
	items, err := repository.ListItems(ctx, transferID)
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
