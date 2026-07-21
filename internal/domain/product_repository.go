package domain

import "context"

// ProductRepository defines the persistence operations required by the product use cases.
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]*Product, error)
	FindByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}
