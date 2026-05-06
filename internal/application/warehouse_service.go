package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// WarehouseService handles warehouse-related business operations.
type WarehouseService struct {
	provider      RepositoryProvider
	notifications NotificationPublisher
}

// NewWarehouseService creates a warehouse service.
func NewWarehouseService(provider RepositoryProvider, publishers ...NotificationPublisher) *WarehouseService {
	return &WarehouseService{
		provider:      provider,
		notifications: notificationPublisherFrom(publishers...),
	}
}

// Create creates an active warehouse with a unique code.
func (service *WarehouseService) Create(ctx context.Context, input CreateWarehouseInput) (*domain.Warehouse, error) {
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	code, err := requireText(input.Code, "code")
	if err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	existing, err := repos.Warehouses.GetByCode(ctx, code)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, conflictError(ErrConflict, "warehouse code already exists")
	}

	warehouse := domain.Warehouse{
		Name:     name,
		Code:     code,
		Location: clean(input.Location),
		IsActive: true,
	}
	if err := repos.Warehouses.Create(ctx, &warehouse); err != nil {
		return nil, err
	}

	return &warehouse, nil
}

// Update updates warehouse details.
func (service *WarehouseService) Update(ctx context.Context, id domain.ID, input UpdateWarehouseInput) (*domain.Warehouse, error) {
	if err := requireID(id, "warehouse id"); err != nil {
		return nil, err
	}
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	code, err := requireText(input.Code, "code")
	if err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	warehouse, err := repos.Warehouses.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing, err := repos.Warehouses.GetByCode(ctx, code)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, conflictError(ErrConflict, "warehouse code already exists")
	}

	wasActive := warehouse.IsActive
	warehouse.Name = name
	warehouse.Code = code
	warehouse.Location = clean(input.Location)
	warehouse.IsActive = input.IsActive

	if err := repos.Warehouses.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	updated, err := repos.Warehouses.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wasActive && !updated.IsActive {
		service.notifications.Publish(ctx, warehouseDeactivatedEvent(*updated))
	}

	return updated, nil
}

// Delete deletes a warehouse when it has no stock or active operations.
func (service *WarehouseService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "warehouse id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := repos.Warehouses.GetByID(ctx, id); err != nil {
		return err
	}

	hasStock, err := repos.Warehouses.HasStock(ctx, id)
	if err != nil {
		return err
	}
	if hasStock {
		return businessRuleError(ErrDeleteRestricted, "warehouse with stock cannot be deleted")
	}

	hasActiveOperations, err := repos.Warehouses.HasActiveOperations(ctx, id)
	if err != nil {
		return err
	}
	if hasActiveOperations {
		return businessRuleError(ErrDeleteRestricted, "warehouse with active operations cannot be deleted")
	}

	movements, err := repos.Movements.ListByWarehouse(ctx, id, ListOptions{Limit: 1})
	if err != nil {
		return err
	}
	if len(movements) > 0 {
		if err := repos.Warehouses.SetActive(ctx, id, false); err != nil {
			return err
		}
		warehouse, err := repos.Warehouses.GetByID(ctx, id)
		if err != nil {
			return err
		}
		service.notifications.Publish(ctx, warehouseDeactivatedEvent(*warehouse))
		return nil
	}

	return repos.Warehouses.Delete(ctx, id)
}

// Activate marks a warehouse as active.
func (service *WarehouseService) Activate(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "warehouse id"); err != nil {
		return err
	}

	return repositories(service.provider).Warehouses.SetActive(ctx, id, true)
}

// Deactivate marks a warehouse as inactive.
func (service *WarehouseService) Deactivate(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "warehouse id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	warehouse, err := repos.Warehouses.GetByID(ctx, id)
	if err != nil {
		return err
	}
	wasActive := warehouse.IsActive

	if err := repos.Warehouses.SetActive(ctx, id, false); err != nil {
		return err
	}
	if wasActive {
		warehouse.IsActive = false
		service.notifications.Publish(ctx, warehouseDeactivatedEvent(*warehouse))
	}

	return nil
}

// GetByID returns a warehouse by identifier.
func (service *WarehouseService) GetByID(ctx context.Context, id domain.ID) (*domain.Warehouse, error) {
	if err := requireID(id, "warehouse id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Warehouses.GetByID(ctx, id)
}

// List returns warehouses matching the filter.
func (service *WarehouseService) List(ctx context.Context, filter WarehouseFilter) ([]domain.Warehouse, error) {
	return repositories(service.provider).Warehouses.List(ctx, filter)
}
