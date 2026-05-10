package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func (handler *Handler) listUsers(ctx *gin.Context) {
	options, err := parseListOptions(ctx)
	if err != nil {
		respondError(ctx, err)
		return
	}

	filter := application.UserFilter{
		Search:      ctx.Query("search"),
		ListOptions: options,
	}
	if role := ctx.Query("role"); role != "" {
		value := parseUserRole(role)
		filter.Role = &value
	}
	isActive, err := parseOptionalBool(ctx, "is_active")
	if err != nil {
		respondError(ctx, err)
		return
	}
	filter.IsActive = isActive

	users, err := handler.services.Users.List(ctx.Request.Context(), filter)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toUserResponses(users))
}

func (handler *Handler) createUser(ctx *gin.Context) {
	var request CreateUserRequest
	if !bindJSON(ctx, &request) {
		return
	}

	user, err := handler.services.Users.Create(ctx.Request.Context(), application.CreateUserInput{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toUserResponse(*user))
}

func (handler *Handler) getUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	user, err := handler.services.Users.GetByID(ctx.Request.Context(), id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toUserResponse(*user))
}

func (handler *Handler) updateUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	var request UpdateUserRequest
	if !bindJSON(ctx, &request) {
		return
	}

	user, err := handler.services.Users.Update(ctx.Request.Context(), id, application.UpdateUserInput{
		Name:     request.Name,
		Email:    request.Email,
		Role:     request.Role,
		IsActive: request.IsActive,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toUserResponse(*user))
}

func (handler *Handler) deleteUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx, "id")
	if err != nil {
		respondError(ctx, err)
		return
	}

	if err := handler.services.Users.Delete(ctx.Request.Context(), id); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func parseUserRole(raw string) domain.UserRole {
	return domain.UserRole(raw)
}
