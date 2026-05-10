package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func (handler *Handler) listSupplies(ctx *gin.Context) {
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	warehouseID, err := parseOptionalID(ctx, "warehouse_id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	supplies, err := handler.services.Supplies.List(ctx.Request.Context(), application.SupplyFilter{
		WarehouseID: warehouseID,
		Status:      parseSupplyStatus(ctx.Query("status")),
		ListOptions: options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toSupplyResponses(supplies))
}

func (handler *Handler) createSupply(ctx *gin.Context) {
	userID, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request SupplyCreateRequest
	if !bindJSON(ctx, &request) {
		return
	}

	supply, err := handler.services.Supplies.Create(ctx.Request.Context(), application.CreateSupplyInput{
		WarehouseID:     request.WarehouseID,
		CreatedByUserID: userID,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toSupplyResponse(*supply))
}

func (handler *Handler) getSupply(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	supply, err := handler.services.Supplies.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toSupplyResponse(*supply))
}

func (handler *Handler) deleteSupply(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Supplies.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) addSupplyItem(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request SupplyItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.Supplies.AddItem(ctx.Request.Context(), id, application.SupplyItemInput{
		ProductID:      request.ProductID,
		Quantity:       request.Quantity,
		UnitPriceCents: request.UnitPriceCents,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toSupplyItemResponse(*item))
}

func (handler *Handler) updateSupplyItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	var request SupplyItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.Supplies.UpdateItem(ctx.Request.Context(), id, itemID, application.SupplyItemInput{
		ProductID:      request.ProductID,
		Quantity:       request.Quantity,
		UnitPriceCents: request.UnitPriceCents,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toSupplyItemResponse(*item))
}

func (handler *Handler) deleteSupplyItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	if err := handler.services.Supplies.DeleteItem(ctx.Request.Context(), id, itemID); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) confirmSupply(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	userID, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}

	supply, err := handler.services.Supplies.Confirm(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toSupplyResponse(*supply))
}

func (handler *Handler) cancelSupply(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	userID, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}

	supply, err := handler.services.Supplies.Cancel(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toSupplyResponse(*supply))
}

func parseOperationItemIDs(ctx *gin.Context) (id domain.ID, itemID domain.ID, ok bool) {
	operationID, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return 0, 0, false
	}
	parsedItemID, err := parseIDParam(ctx, "itemID")
	if err != nil {
		respondError(ctx, err)
		return 0, 0, false
	}

	return operationID, parsedItemID, true
}
