package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func (handler *Handler) listWarehouses(ctx *gin.Context) {
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	isActive, err := parseOptionalBool(ctx, "is_active")
	if err != nil {
		respondError(ctx, err)
		return
	}

	warehouses, err := handler.services.Warehouses.List(ctx.Request.Context(), application.WarehouseFilter{
		IsActive:    isActive,
		Code:        ctx.Query("code"),
		Search:      ctx.Query("search"),
		ListOptions: options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toWarehouseResponses(warehouses))
}

func (handler *Handler) createWarehouse(ctx *gin.Context) {
	var request WarehouseRequest
	if !bindJSON(ctx, &request) {
		return
	}

	warehouse, err := handler.services.Warehouses.Create(ctx.Request.Context(), application.CreateWarehouseInput{
		Name:     request.Name,
		Code:     request.Code,
		Location: request.Location,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toWarehouseResponse(*warehouse))
}

func (handler *Handler) getWarehouse(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	warehouse, err := handler.services.Warehouses.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toWarehouseResponse(*warehouse))
}

func (handler *Handler) updateWarehouse(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request WarehouseRequest
	if !bindJSON(ctx, &request) {
		return
	}

	warehouse, err := handler.services.Warehouses.Update(ctx.Request.Context(), id, application.UpdateWarehouseInput{
		Name:     request.Name,
		Code:     request.Code,
		Location: request.Location,
		IsActive: request.IsActive,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toWarehouseResponse(*warehouse))
}

func (handler *Handler) deleteWarehouse(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Warehouses.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) activateWarehouse(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Warehouses.Activate(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) deactivateWarehouse(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Warehouses.Deactivate(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) listCategories(ctx *gin.Context) {
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	isActive, err := parseOptionalBool(ctx, "is_active")
	if err != nil {
		respondError(ctx, err)
		return
	}

	categories, err := handler.services.Categories.List(ctx.Request.Context(), application.CategoryFilter{
		IsActive:    isActive,
		Search:      ctx.Query("search"),
		ListOptions: options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toCategoryResponses(categories))
}

func (handler *Handler) createCategory(ctx *gin.Context) {
	var request CategoryRequest
	if !bindJSON(ctx, &request) {
		return
	}

	category, err := handler.services.Categories.Create(ctx.Request.Context(), application.CreateCategoryInput{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toCategoryResponse(*category))
}

func (handler *Handler) getCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	category, err := handler.services.Categories.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCategoryResponse(*category))
}

func (handler *Handler) updateCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request CategoryRequest
	if !bindJSON(ctx, &request) {
		return
	}

	category, err := handler.services.Categories.Update(ctx.Request.Context(), id, application.UpdateCategoryInput{
		Name:        request.Name,
		Description: request.Description,
		IsActive:    request.IsActive,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCategoryResponse(*category))
}

func (handler *Handler) deleteCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Categories.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) listProducts(ctx *gin.Context) {
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	categoryID, err := parseOptionalID(ctx, "category_id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	isActive, err := parseOptionalBool(ctx, "is_active")
	if err != nil {
		respondError(ctx, err)
		return
	}

	products, err := handler.services.Products.List(ctx.Request.Context(), application.ProductFilter{
		CategoryID:  categoryID,
		IsActive:    isActive,
		SKU:         ctx.Query("sku"),
		Name:        ctx.Query("name"),
		Search:      ctx.Query("search"),
		ListOptions: options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toProductResponses(products))
}

func (handler *Handler) createProduct(ctx *gin.Context) {
	var request ProductRequest
	if !bindJSON(ctx, &request) {
		return
	}

	product, err := handler.services.Products.Create(ctx.Request.Context(), application.CreateProductInput{
		Name:              request.Name,
		SKU:               request.SKU,
		CategoryID:        request.CategoryID,
		UnitOfMeasure:     request.UnitOfMeasure,
		MinStockThreshold: request.MinStockThreshold,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toProductResponse(*product))
}

func (handler *Handler) getProduct(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	product, err := handler.services.Products.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toProductResponse(*product))
}

func (handler *Handler) updateProduct(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request ProductRequest
	if !bindJSON(ctx, &request) {
		return
	}

	product, err := handler.services.Products.Update(ctx.Request.Context(), id, application.UpdateProductInput{
		Name:              request.Name,
		SKU:               request.SKU,
		CategoryID:        request.CategoryID,
		UnitOfMeasure:     request.UnitOfMeasure,
		MinStockThreshold: request.MinStockThreshold,
		IsActive:          request.IsActive,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toProductResponse(*product))
}

func (handler *Handler) deleteProduct(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Products.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) activateProduct(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Products.Activate(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) deactivateProduct(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Products.Deactivate(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func parseSupplyStatus(raw string) *domain.SupplyStatus {
	if raw == "" {
		return nil
	}
	status := domain.SupplyStatus(raw)
	return &status
}

func parseOutboundStatus(raw string) *domain.OutboundRequestStatus {
	if raw == "" {
		return nil
	}
	status := domain.OutboundRequestStatus(raw)
	return &status
}

func parseTransferStatus(raw string) *domain.TransferStatus {
	if raw == "" {
		return nil
	}
	status := domain.TransferStatus(raw)
	return &status
}
