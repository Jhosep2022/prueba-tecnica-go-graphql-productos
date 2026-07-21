package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

// ProductUseCase contains the business rules for managing products.
type ProductUseCase struct {
	repository domain.ProductRepository
}

// NewProductUseCase creates a ProductUseCase with its required repository.
func NewProductUseCase(repository domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		repository: repository,
	}
}

// CreateProduct validates and creates a new product.
func (uc *ProductUseCase) CreateProduct(
	ctx context.Context,
	name string,
	price float64,
	stock int,
) (*domain.Product, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, &domain.ValidationError{
			Field:   "name",
			Message: "must not be empty",
		}
	}

	if price <= 0 {
		return nil, &domain.ValidationError{
			Field:   "price",
			Message: "must be greater than zero",
		}
	}

	if stock < 0 {
		return nil, &domain.ValidationError{
			Field:   "stock",
			Message: "must not be negative",
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("generate product id: %w", err)
	}

	product := &domain.Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.repository.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

// ListProducts returns all existing products.
func (uc *ProductUseCase) ListProducts(
	ctx context.Context,
) ([]*domain.Product, error) {
	products, err := uc.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return products, nil
}

// GetProductByID returns a product using its identifier.
func (uc *ProductUseCase) GetProductByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, &domain.ValidationError{
			Field:   "id",
			Message: "must not be empty",
		}
	}

	product, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates the name or price of an existing product.
func (uc *ProductUseCase) UpdateProduct(
	ctx context.Context,
	id string,
	name *string,
	price *float64,
) (*domain.Product, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, &domain.ValidationError{
			Field:   "id",
			Message: "must not be empty",
		}
	}

	if name == nil && price == nil {
		return nil, &domain.ValidationError{
			Field:   "input",
			Message: "name or price must be provided",
		}
	}

	product, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product to update: %w", err)
	}

	if name != nil {
		trimmedName := strings.TrimSpace(*name)

		if trimmedName == "" {
			return nil, &domain.ValidationError{
				Field:   "name",
				Message: "must not be empty",
			}
		}

		product.Name = trimmedName
	}

	if price != nil {
		if *price <= 0 {
			return nil, &domain.ValidationError{
				Field:   "price",
				Message: "must be greater than zero",
			}
		}

		product.Price = *price
	}

	if err := uc.repository.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return product, nil
}

// DeleteProduct removes a product using its identifier.
func (uc *ProductUseCase) DeleteProduct(
	ctx context.Context,
	id string,
) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return &domain.ValidationError{
			Field:   "id",
			Message: "must not be empty",
		}
	}

	if err := uc.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}

// generateID creates a UUID-compatible identifier using the standard library.
func generateID() (string, error) {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		value[0:4],
		value[4:6],
		value[6:8],
		value[8:10],
		value[10:16],
	), nil
}
