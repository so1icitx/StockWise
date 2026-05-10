package httpapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func parseIDParam(ctx *gin.Context, name string) (domain.ID, error) {
	raw := ctx.Param(name)
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return 0, validationHTTPError(name + " must be a positive integer")
	}

	return domain.ID(value), nil
}

func parseOptionalBool(ctx *gin.Context, name string) (*bool, error) {
	raw := ctx.Query(name)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, validationHTTPError(name + " must be true or false")
	}

	return &value, nil
}

func parseOptionalID(ctx *gin.Context, name string) (*domain.ID, error) {
	raw := ctx.Query(name)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return nil, validationHTTPError(name + " must be a positive integer")
	}

	id := domain.ID(value)
	return &id, nil
}

func parseListOptions(ctx *gin.Context) (application.ListOptions, error) {
	limit, err := parseOptionalInt(ctx, "limit")
	if err != nil {
		return application.ListOptions{}, err
	}
	offset, err := parseOptionalInt(ctx, "offset")
	if err != nil {
		return application.ListOptions{}, err
	}

	return application.ListOptions{Limit: limit, Offset: offset}, nil
}

func parseOptionalInt(ctx *gin.Context, name string) (int, error) {
	raw := ctx.Query(name)
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, validationHTTPError(name + " must be a non-negative integer")
	}

	return value, nil
}
