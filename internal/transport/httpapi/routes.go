package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
)

// RegisterRoutes registers StockWise REST API routes.
func RegisterRoutes(api *gin.RouterGroup, services application.Services) {
	handler := &Handler{services: services}

	users := api.Group("/users")
	users.GET("", handler.listUsers)
	users.POST("", handler.createUser)
	users.GET("/:id", handler.getUser)
	users.PUT("/:id", handler.updateUser)
	users.DELETE("/:id", handler.deleteUser)

	warehouses := api.Group("/warehouses")
	warehouses.GET("", handler.listWarehouses)
	warehouses.POST("", handler.createWarehouse)
	warehouses.GET("/:id", handler.getWarehouse)
	warehouses.PUT("/:id", handler.updateWarehouse)
	warehouses.DELETE("/:id", handler.deleteWarehouse)
	warehouses.PATCH("/:id/activate", handler.activateWarehouse)
	warehouses.PATCH("/:id/deactivate", handler.deactivateWarehouse)
	warehouses.GET("/:id/stock", handler.getWarehouseStock)
	warehouses.GET("/:id/movements", handler.getWarehouseMovements)

	categories := api.Group("/categories")
	categories.GET("", handler.listCategories)
	categories.POST("", handler.createCategory)
	categories.GET("/:id", handler.getCategory)
	categories.PUT("/:id", handler.updateCategory)
	categories.DELETE("/:id", handler.deleteCategory)

	products := api.Group("/products")
	products.GET("", handler.listProducts)
	products.POST("", handler.createProduct)
	products.GET("/low-stock", handler.getLowStockProducts)
	products.GET("/:id", handler.getProduct)
	products.PUT("/:id", handler.updateProduct)
	products.DELETE("/:id", handler.deleteProduct)
	products.PATCH("/:id/activate", handler.activateProduct)
	products.PATCH("/:id/deactivate", handler.deactivateProduct)
	products.GET("/:id/total-stock", handler.getProductTotalStock)
	products.GET("/:id/stock", handler.getProductStock)
	products.GET("/:id/movements", handler.getProductMovements)

	supplies := api.Group("/supplies")
	supplies.GET("", handler.listSupplies)
	supplies.POST("", handler.createSupply)
	supplies.GET("/:id", handler.getSupply)
	supplies.DELETE("/:id", handler.deleteSupply)
	supplies.POST("/:id/items", handler.addSupplyItem)
	supplies.PUT("/:id/items/:itemID", handler.updateSupplyItem)
	supplies.DELETE("/:id/items/:itemID", handler.deleteSupplyItem)
	supplies.POST("/:id/confirm", handler.confirmSupply)
	supplies.POST("/:id/cancel", handler.cancelSupply)

	outboundRequests := api.Group("/outbound-requests")
	outboundRequests.GET("", handler.listOutboundRequests)
	outboundRequests.POST("", handler.createOutboundRequest)
	outboundRequests.GET("/:id", handler.getOutboundRequest)
	outboundRequests.DELETE("/:id", handler.deleteOutboundRequest)
	outboundRequests.POST("/:id/items", handler.addOutboundRequestItem)
	outboundRequests.PUT("/:id/items/:itemID", handler.updateOutboundRequestItem)
	outboundRequests.DELETE("/:id/items/:itemID", handler.deleteOutboundRequestItem)
	outboundRequests.POST("/:id/approve", handler.approveOutboundRequest)
	outboundRequests.POST("/:id/reject", handler.rejectOutboundRequest)
	outboundRequests.POST("/:id/cancel", handler.cancelOutboundRequest)
	outboundRequests.POST("/:id/execute", handler.executeOutboundRequest)

	transfers := api.Group("/transfers")
	transfers.GET("", handler.listTransfers)
	transfers.POST("", handler.createTransfer)
	transfers.GET("/:id", handler.getTransfer)
	transfers.DELETE("/:id", handler.deleteTransfer)
	transfers.POST("/:id/items", handler.addTransferItem)
	transfers.PUT("/:id/items/:itemID", handler.updateTransferItem)
	transfers.DELETE("/:id/items/:itemID", handler.deleteTransferItem)
	transfers.POST("/:id/confirm", handler.confirmTransfer)
	transfers.POST("/:id/cancel", handler.cancelTransfer)
}

// Handler contains REST route handlers backed by application services.
type Handler struct {
	services application.Services
}
