package graphql

import "github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/usecase"

// Resolver contains the dependencies for required by GraphQL resolvers.
type Resolver struct {
	productUseCase *usecase.ProductUseCase
}

// NewResolver creates a new Resolver with its required dependencies.
func NewResolver(productUseCase *usecase.ProductUseCase) *Resolver {
	return &Resolver{
		productUseCase: productUseCase,
	}
}
