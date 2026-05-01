package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/config"
	graphqlapi "github.com/so1icitx/StockWise/internal/transport/graphql"
	websocketapi "github.com/so1icitx/StockWise/internal/transport/websocket"
)

// NewRouter builds the Gin router for the StockWise REST API.
func NewRouter(cfg config.Config, services application.Services, notificationHubs ...*websocketapi.Hub) *gin.Engine {
	if cfg.AppEnv == "test" {
		gin.SetMode(gin.TestMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), ActorMiddleware())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
			"env":    cfg.AppEnv,
		})
	})

	api := router.Group("/api/v1")
	api.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
			"env":    cfg.AppEnv,
		})
	})

	RegisterRoutes(api, services)
	graphqlapi.RegisterRoutes(router, services)
	if len(notificationHubs) > 0 {
		websocketapi.RegisterRoutes(router, notificationHubs[0])
	}
	RegisterSwagger(router)

	return router
}
