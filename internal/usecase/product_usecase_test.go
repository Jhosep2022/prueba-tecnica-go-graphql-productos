package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

// fakeProductRepository allows the use cases to be tested without real storage.
type fakeProductRepository struct {
	createFn   func(context.Context, *domain.Product) error
	findAllFn  func(context.Context) ([]*domain.Product, error)
	findByIDFn func(context.Context, string) (*domain.Product, error)
	updateFn   func(context.Context, *domain.Product) error
	deleteFn   func(context.Context, string) error
}

func (r *fakeProductRepository) Create(
	ctx context.Context,
	product *domain.Product,
) error {
	if r.createFn == nil {
		return nil
	}

	return r.createFn(ctx, product)
}

func (r *fakeProductRepository) FindAll(
	ctx context.Context,
) ([]*domain.Product, error) {
	if r.findAllFn == nil {
		return []*domain.Product{}, nil
	}

	return r.findAllFn(ctx)
}

func (r *fakeProductRepository) FindByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {
	if r.findByIDFn == nil {
		return nil, nil
	}

	return r.findByIDFn(ctx, id)
}

func (r *fakeProductRepository) Update(
	ctx context.Context,
	product *domain.Product,
) error {
	if r.updateFn == nil {
		return nil
	}

	return r.updateFn(ctx, product)
}

func (r *fakeProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	if r.deleteFn == nil {
		return nil
	}

	return r.deleteFn(ctx, id)
}

func TestCreateProductSuccess(t *testing.T) {
	var savedProduct *domain.Product

	repository := &fakeProductRepository{
		createFn: func(_ context.Context, product *domain.Product) error {
			savedProduct = product
			return nil
		},
	}

	productUseCase := NewProductUseCase(repository)

	product, err := productUseCase.CreateProduct(
		context.Background(),
		"  Keyboard  ",
		120.50,
		10,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID == "" {
		t.Error("expected product ID to be generated")
	}

	if product.Name != "Keyboard" {
		t.Errorf("expected trimmed name Keyboard, got %q", product.Name)
	}

	if product.Price != 120.50 {
		t.Errorf("expected price 120.50, got %v", product.Price)
	}

	if product.Stock != 10 {
		t.Errorf("expected stock 10, got %d", product.Stock)
	}

	if product.CreatedAt.IsZero() {
		t.Error("expected creation date to be assigned")
	}

	if savedProduct != product {
		t.Error("expected product to be sent to repository")
	}
}

func TestCreateProductValidation(t *testing.T) {
	tests := []struct {
		name          string
		productName   string
		price         float64
		stock         int
		expectedField string
	}{
		{
			name:          "empty name",
			productName:   "   ",
			price:         10,
			stock:         1,
			expectedField: "name",
		},
		{
			name:          "invalid price",
			productName:   "Monitor",
			price:         0,
			stock:         1,
			expectedField: "price",
		},
		{
			name:          "negative stock",
			productName:   "Monitor",
			price:         10,
			stock:         -1,
			expectedField: "stock",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryCalled := false

			repository := &fakeProductRepository{
				createFn: func(_ context.Context, _ *domain.Product) error {
					repositoryCalled = true
					return nil
				},
			}

			productUseCase := NewProductUseCase(repository)

			product, err := productUseCase.CreateProduct(
				context.Background(),
				test.productName,
				test.price,
				test.stock,
			)

			if err == nil {
				t.Fatal("expected validation error")
			}

			if product != nil {
				t.Error("expected product to be nil")
			}

			var validationError *domain.ValidationError
			if !errors.As(err, &validationError) {
				t.Fatalf("expected ValidationError, got %T", err)
			}

			if validationError.Field != test.expectedField {
				t.Errorf(
					"expected validation field %q, got %q",
					test.expectedField,
					validationError.Field,
				)
			}

			if repositoryCalled {
				t.Error("repository must not be called when validation fails")
			}
		})
	}
}

func TestUpdateProductSuccess(t *testing.T) {
	createdAt := time.Now().UTC()

	storedProduct := &domain.Product{
		ID:        "product-1",
		Name:      "Old name",
		Price:     50,
		Stock:     8,
		CreatedAt: createdAt,
	}

	repository := &fakeProductRepository{
		findByIDFn: func(_ context.Context, id string) (*domain.Product, error) {
			if id != "product-1" {
				t.Errorf("expected ID product-1, got %q", id)
			}

			return storedProduct, nil
		},
		updateFn: func(_ context.Context, product *domain.Product) error {
			if product.Name != "New name" {
				t.Errorf("expected updated name, got %q", product.Name)
			}

			if product.Price != 75 {
				t.Errorf("expected updated price 75, got %v", product.Price)
			}

			return nil
		},
	}

	productUseCase := NewProductUseCase(repository)

	newName := "  New name  "
	newPrice := 75.0

	product, err := productUseCase.UpdateProduct(
		context.Background(),
		"product-1",
		&newName,
		&newPrice,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Stock != 8 {
		t.Errorf("expected stock to remain unchanged, got %d", product.Stock)
	}

	if !product.CreatedAt.Equal(createdAt) {
		t.Error("expected creation date to remain unchanged")
	}
}

func TestGetProductByIDNotFound(t *testing.T) {
	repository := &fakeProductRepository{
		findByIDFn: func(_ context.Context, id string) (*domain.Product, error) {
			return nil, &domain.ProductNotFoundError{ID: id}
		},
	}

	productUseCase := NewProductUseCase(repository)

	product, err := productUseCase.GetProductByID(
		context.Background(),
		"missing-product",
	)

	if err == nil {
		t.Fatal("expected product not found error")
	}

	if product != nil {
		t.Error("expected product to be nil")
	}

	var notFoundError *domain.ProductNotFoundError
	if !errors.As(err, &notFoundError) {
		t.Fatalf("expected ProductNotFoundError, got %T", err)
	}

	if notFoundError.ID != "missing-product" {
		t.Errorf("expected missing-product ID, got %q", notFoundError.ID)
	}
}
