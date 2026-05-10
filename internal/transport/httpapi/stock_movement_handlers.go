package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) getWarehouseStock(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	stockItems, err := handler.services.Stock.GetByWarehouse(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toStockItemResponses(stockItems))
}

func (handler *Handler) getProductStock(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	stockItems, err := handler.services.Stock.GetByProductAcrossWarehouses(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toStockItemResponses(stockItems))
}

func (handler *Handler) getProductTotalStock(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	quantity, err := handler.services.Stock.GetTotalQuantityForProduct(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, TotalStockResponse{ProductID: id, Quantity: quantity})
}

func (handler *Handler) getLowStockProducts(ctx *gin.Context) {
	statuses, err := handler.services.Stock.GetLowStock(ctx.Request.Context())
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toStockStatusResponses(statuses))
}

func (handler *Handler) getProductMovements(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}

	movements, err := handler.services.Movements.ListByProduct(ctx.Request.Context(), id, options)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toMovementResponses(movements))
}

func (handler *Handler) getWarehouseMovements(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}

	movements, err := handler.services.Movements.ListByWarehouse(ctx.Request.Context(), id, options)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toMovementResponses(movements))
}
