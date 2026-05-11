package graphqlapi

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/transport/graphql/generated"
)

// RegisterRoutes registers GraphQL endpoint and playground routes.
func RegisterRoutes(router *gin.Engine, services application.Services) {
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: NewResolver(services),
	}))

	router.Any("/graphql", gin.WrapH(server))
	router.GET("/graphql/playground", gin.WrapH(playground.Handler("StockWise GraphQL", "/graphql")))
}
