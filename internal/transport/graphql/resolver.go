package graphqlapi

import "github.com/so1icitx/StockWise/internal/application"

// Resolver provides gqlgen resolvers backed by application services.
type Resolver struct {
	services application.Services
}

// NewResolver creates a GraphQL resolver root.
func NewResolver(services application.Services) *Resolver {
	return &Resolver{services: services}
}
