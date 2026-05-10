package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func (handler *Handler) listOutboundRequests(ctx *gin.Context) {
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

	requests, err := handler.services.OutboundRequests.List(ctx.Request.Context(), application.OutboundRequestFilter{
		WarehouseID: warehouseID,
		Status:      parseOutboundStatus(ctx.Query("status")),
		ListOptions: options,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toOutboundRequestResponses(requests))
}

func (handler *Handler) createOutboundRequest(ctx *gin.Context) {
	userID, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request OutboundRequestCreateRequest
	if !bindJSON(ctx, &request) {
		return
	}

	outboundRequest, err := handler.services.OutboundRequests.Create(ctx.Request.Context(), application.CreateOutboundRequestInput{
		WarehouseID:     request.WarehouseID,
		CreatedByUserID: userID,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toOutboundRequestResponse(*outboundRequest))
}

func (handler *Handler) getOutboundRequest(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	request, err := handler.services.OutboundRequests.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toOutboundRequestResponse(*request))
}

func (handler *Handler) deleteOutboundRequest(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	if err := handler.services.OutboundRequests.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) addOutboundRequestItem(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}
	var request OutboundRequestItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.OutboundRequests.AddItem(ctx.Request.Context(), id, application.OutboundRequestItemInput{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toOutboundRequestItemResponse(*item))
}

func (handler *Handler) updateOutboundRequestItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	var request OutboundRequestItemRequest
	if !bindJSON(ctx, &request) {
		return
	}

	item, err := handler.services.OutboundRequests.UpdateItem(ctx.Request.Context(), id, itemID, application.OutboundRequestItemInput{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toOutboundRequestItemResponse(*item))
}

func (handler *Handler) deleteOutboundRequestItem(ctx *gin.Context) {
	id, itemID, ok := parseOperationItemIDs(ctx)
	if !ok {
		return
	}
	if err := handler.services.OutboundRequests.DeleteItem(ctx.Request.Context(), id, itemID); err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) approveOutboundRequest(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	request, err := handler.services.OutboundRequests.Approve(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toOutboundRequestResponse(*request))
}

func (handler *Handler) rejectOutboundRequest(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	request, err := handler.services.OutboundRequests.Reject(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toOutboundRequestResponse(*request))
}

func (handler *Handler) cancelOutboundRequest(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	request, err := handler.services.OutboundRequests.Cancel(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toOutboundRequestResponse(*request))
}

func (handler *Handler) executeOutboundRequest(ctx *gin.Context) {
	id, userID, ok := parseOperationAndActor(ctx)
	if !ok {
		return
	}
	request, err := handler.services.OutboundRequests.Execute(ctx.Request.Context(), id, userID)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toOutboundRequestResponse(*request))
}

func parseOperationAndActor(ctx *gin.Context) (id domain.ID, userID domain.ID, ok bool) {
	operationID, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return 0, 0, false
	}
	actor, err := actorID(ctx)
	if err != nil {
		respondError(ctx, err)
		return 0, 0, false
	}

	return operationID, actor, true
}
