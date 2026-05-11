package graphqlapi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func parseID(raw string, field string) (domain.ID, error) {
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return 0, fmt.Errorf("%s must be a positive ID", field)
	}

	return domain.ID(value), nil
}

func parseOptionalID(raw *string, field string) (*domain.ID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}

	value, err := parseID(*raw, field)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func idString(id domain.ID) string {
	return strconv.FormatUint(uint64(id), 10)
}

func idStringPtr(id *domain.ID) *string {
	if id == nil {
		return nil
	}

	value := idString(*id)
	return &value
}

func timeString(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func timeStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := timeString(*value)
	return &formatted
}

func listOptions(limit *int, offset *int) application.ListOptions {
	options := application.ListOptions{}
	if limit != nil {
		options.Limit = *limit
	}
	if offset != nil {
		options.Offset = *offset
	}

	return options
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func optionalStringWithFallback(value *string, fallback string) string {
	if value == nil {
		return fallback
	}

	return *value
}

func optionalBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}

	return *value
}

func parseOperationActorIDs(operationID string, actorID string, operationField string) (domain.ID, domain.ID, error) {
	operation, err := parseID(operationID, operationField)
	if err != nil {
		return 0, 0, err
	}
	actor, err := parseID(actorID, "user id")
	if err != nil {
		return 0, 0, err
	}

	return operation, actor, nil
}

func parseOperationItemIDs(operationID string, itemID string, operationField string, itemField string) (domain.ID, domain.ID, error) {
	operation, err := parseID(operationID, operationField)
	if err != nil {
		return 0, 0, err
	}
	item, err := parseID(itemID, itemField)
	if err != nil {
		return 0, 0, err
	}

	return operation, item, nil
}

func parseOperationAndProductIDs(operationID string, productID string, operationField string) (domain.ID, domain.ID, error) {
	operation, err := parseID(operationID, operationField)
	if err != nil {
		return 0, 0, err
	}
	product, err := parseID(productID, "product id")
	if err != nil {
		return 0, 0, err
	}

	return operation, product, nil
}

func parseOperationItemAndProductIDs(operationID string, itemID string, productID string, operationField string, itemField string) (domain.ID, domain.ID, domain.ID, error) {
	operation, item, err := parseOperationItemIDs(operationID, itemID, operationField, itemField)
	if err != nil {
		return 0, 0, 0, err
	}
	product, err := parseID(productID, "product id")
	if err != nil {
		return 0, 0, 0, err
	}

	return operation, item, product, nil
}
