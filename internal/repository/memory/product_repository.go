package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

// ProductRepository is an in-memory implementation of the ProductRepository interface.
type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]*domain.Product
}

// NewProductRepository creates an empty in-memory product repository.
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]*domain.Product),
	}
}

// Create stores a new product.
func (r *ProductRepository) Create(
	ctx context.Context,
	product *domain.Product,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if product == nil {
		return &domain.ValidationError{
			Field:   "product",
			Message: "must not be nil",
		}
	}

	if product.ID == "" {
		return &domain.ValidationError{
			Field:   "id",
			Message: "must not be empty",
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return &domain.ProductAlreadyExistsError{
			ID: product.ID,
		}
	}

	r.products[product.ID] = cloneProduct(product)
	return nil
}

// FindAll returns all products ordered by creation date.
func (r *ProductRepository) FindAll(
	ctx context.Context,
) ([]*domain.Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*domain.Product, 0, len(r.products))

	for _, product := range r.products {
		products = append(products, cloneProduct(product))
	}

	sort.Slice(products, func(i, j int) bool {
		if products[i].CreatedAt.Equal(products[j].CreatedAt) {
			return products[i].ID < products[j].ID
		}

		return products[i].CreatedAt.Before(products[j].CreatedAt)
	})

	return products, nil
}

// FindByID returns a product using its identifier.
func (r *ProductRepository) FindByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, &domain.ProductNotFoundError{ID: id}
	}

	return cloneProduct(product), nil
}

// Update replaces an existing product.
func (r *ProductRepository) Update(
	ctx context.Context,
	product *domain.Product,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if product == nil {
		return &domain.ValidationError{
			Field:   "product",
			Message: "must not be nil",
		}
	}

	if product.ID == "" {
		return &domain.ValidationError{
			Field:   "id",
			Message: "must not be empty",
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return &domain.ProductNotFoundError{
			ID: product.ID,
		}
	}

	r.products[product.ID] = cloneProduct(product)

	return nil
}

// Delete removes a product using its identifier.
func (r *ProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return &domain.ProductNotFoundError{ID: id}
	}

	delete(r.products, id)

	return nil
}

// cloneProduct creates a copy so callers cannot mutate stored data directly.
func cloneProduct(product *domain.Product) *domain.Product {
	if product == nil {
		return nil
	}

	copy := *product

	return &copy
}
