package application

import "errors"

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ErrValidation is returned when input data is invalid.
var ErrValidation = errors.New("validation failed")

// ErrConflict is returned when a request conflicts with existing state.
var ErrConflict = errors.New("conflict")

// ErrBusinessRule is returned when a domain rule prevents an operation.
var ErrBusinessRule = errors.New("business rule violation")

// ErrDuplicateSKU is returned when a product SKU is already used.
var ErrDuplicateSKU = errors.New("duplicate product sku")

// ErrInactiveWarehouse is returned when an inactive warehouse is used for an operation.
var ErrInactiveWarehouse = errors.New("inactive warehouse")

// ErrInactiveProduct is returned when an inactive product is used for an operation.
var ErrInactiveProduct = errors.New("inactive product")

// ErrInactiveCategory is returned when an inactive category is used for an active product.
var ErrInactiveCategory = errors.New("inactive category")

// ErrInsufficientStock is returned when available stock cannot satisfy an operation.
var ErrInsufficientStock = errors.New("insufficient stock")

// ErrOperationLocked is returned when a completed or cancelled operation is edited.
var ErrOperationLocked = errors.New("operation cannot be edited")

// ErrDeleteRestricted is returned when deletion would break stock or movement history.
var ErrDeleteRestricted = errors.New("delete restricted")

// RuleError wraps a business, validation, or conflict error with a human-readable message.
type RuleError struct {
	Kind    error
	Message string
}

// Error returns the human-readable rule error message.
func (err RuleError) Error() string {
	return err.Message
}

// Unwrap returns the underlying error kind for errors.Is checks.
func (err RuleError) Unwrap() error {
	return err.Kind
}

func validationError(message string) error {
	return RuleError{Kind: ErrValidation, Message: message}
}

func conflictError(kind error, message string) error {
	return RuleError{Kind: kind, Message: message}
}

func businessRuleError(kind error, message string) error {
	return RuleError{Kind: kind, Message: message}
}
