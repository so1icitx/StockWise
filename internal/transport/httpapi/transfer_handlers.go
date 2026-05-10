package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
)

func (handler *Handler) listTransfers(ctx *gin.Context) {
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
	sourceWarehouseID, err := parseOptionalID(ctx, "source_warehouse_id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	targetWarehouseID, err := parseOptionalID(ctx, "target_warehouse_id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	transfers, err := handler.services.Transfers.List(ctx.Request.Context(), application.TransferFilter{
		WarehouseID:       warehouseID,
		SourceWarehouseID: sourceWarehouseID,
		TargetWarehouseID: targetWarehouseID,
		Status:            parseTransferStatus(ctx.Query("status")),
		ListOptions:       options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toTransferResponses(transfers))
}

func (handler *Handler) createTransfer(ctx *gin.Context) {
	userID, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request TransferCreateRequest
	if !bindJSON(ctx, &request) {
		return
	}

	transfer, err := handler.services.Transfers.Create(ctx.Request.Context(), application.CreateTransferInput{
		SourceWarehouseID: request.SourceWarehouseID,
		TargetWarehouseID: request.TargetWarehouseID,
		CreatedByUserID:   userID,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toTransferResponse(*transfer))
}

func (handler *Handler) getTransfer(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	transfer, err := handler.services.Transfers.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransferResponse(*transfer))
}

func (handler *Handler) deleteTransfer(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.Transfers.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) addTransferItem(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request TransferItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.Transfers.AddItem(ctx.Request.Context(), id, application.TransferItemInput{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toTransferItemResponse(*item))
}

func (handler *Handler) updateTransferItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	var request TransferItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.Transfers.UpdateItem(ctx.Request.Context(), id, itemID, application.TransferItemInput{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransferItemResponse(*item))
}

func (handler *Handler) deleteTransferItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	if err := handler.services.Transfers.DeleteItem(ctx.Request.Context(), id, itemID); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) confirmTransfer(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	transfer, err := handler.services.Transfers.Confirm(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransferResponse(*transfer))
}

func (handler *Handler) cancelTransfer(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	transfer, err := handler.services.Transfers.Cancel(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransferResponse(*transfer))
}
