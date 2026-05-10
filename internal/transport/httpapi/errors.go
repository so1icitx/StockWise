package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
)

// ErrorResponse is the standard REST error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"validation failed"`
}

type httpValidationError struct {
	message string
}

func (err httpValidationError) Error() string {
	return err.message
}

func validationHTTPError(message string) error {
	return httpValidationError{message: message}
}

func respondError(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.As(err, new(httpValidationError)):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, application.ErrValidation):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, application.ErrNotFound):
		status = http.StatusNotFound
		message = "resource not found"
	case errors.Is(err, application.ErrConflict),
		errors.Is(err, application.ErrBusinessRule),
		errors.Is(err, application.ErrDuplicateSKU),
		errors.Is(err, application.ErrInactiveWarehouse),
		errors.Is(err, application.ErrInactiveProduct),
		errors.Is(err, application.ErrInactiveCategory),
		errors.Is(err, application.ErrInsufficientStock),
		errors.Is(err, application.ErrOperationLocked),
		errors.Is(err, application.ErrDeleteRestricted):
		status = http.StatusConflict
		message = err.Error()
	}

	ctx.JSON(status, ErrorResponse{Error: message})
}

func bindJSON(ctx *gin.Context, target any) bool {
	if err := ctx.ShouldBindJSON(target); err != nil {
		respondError(ctx, validationHTTPError(err.Error()))
		return false
	}

	return true
}
