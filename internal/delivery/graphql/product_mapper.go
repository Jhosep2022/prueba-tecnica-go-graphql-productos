package graphql

import (
	"time"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/graph/model"
	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

func toProductModel(product *domain.Product) *model.Product {
	if product == nil {
		return nil
	}

	return &model.Product{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
	}
}

func toProductModels(products []*domain.Product) []*model.Product {
	result := make([]*model.Product, 0, len(products))

	for _, product := range products {
		result = append(result, toProductModel(product))
	}

	return result
}
