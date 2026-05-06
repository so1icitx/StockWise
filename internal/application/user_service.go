package application

import (
	"context"
	"strings"

	"github.com/so1icitx/StockWise/internal/domain"
)

// UserService handles user-related business operations.
type UserService struct {
	provider RepositoryProvider
}

// NewUserService creates a user service.
func NewUserService(provider RepositoryProvider) *UserService {
	return &UserService{provider: provider}
}

// Create creates a new active user.
func (service *UserService) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if err := validateUserRole(input.Role); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	existing, err := repos.Users.GetByEmail(ctx, email)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, conflictError(ErrConflict, "user email already exists")
	}

	user := domain.User{
		Name:     name,
		Email:    email,
		Role:     input.Role,
		IsActive: true,
	}
	if err := repos.Users.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user.
func (service *UserService) Update(ctx context.Context, id domain.ID, input UpdateUserInput) (*domain.User, error) {
	if err := requireID(id, "user id"); err != nil {
		return nil, err
	}
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if err := validateUserRole(input.Role); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	user, err := repos.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing, err := repos.Users.GetByEmail(ctx, email)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, conflictError(ErrConflict, "user email already exists")
	}

	user.Name = name
	user.Email = email
	user.Role = input.Role
	user.IsActive = input.IsActive

	if err := repos.Users.Update(ctx, user); err != nil {
		return nil, err
	}

	return repos.Users.GetByID(ctx, id)
}

// Delete removes a user.
func (service *UserService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "user id"); err != nil {
		return err
	}

	return repositories(service.provider).Users.Delete(ctx, id)
}

// GetByID returns a user by identifier.
func (service *UserService) GetByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	if err := requireID(id, "user id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Users.GetByID(ctx, id)
}

// GetByEmail returns a user by email address.
func (service *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	return repositories(service.provider).Users.GetByEmail(ctx, normalized)
}

// List returns users matching the filter.
func (service *UserService) List(ctx context.Context, filter UserFilter) ([]domain.User, error) {
	return repositories(service.provider).Users.List(ctx, filter)
}

func validateUserRole(role domain.UserRole) error {
	switch role {
	case domain.UserRoleAdmin, domain.UserRoleManager, domain.UserRoleOperator:
		return nil
	default:
		return validationError("user role is invalid")
	}
}

func normalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", validationError("email is required")
	}
	if !strings.Contains(normalized, "@") {
		return "", validationError("email must contain @")
	}

	return normalized, nil
}
