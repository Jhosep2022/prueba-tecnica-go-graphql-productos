package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

func TestProductRepositoryCreateAndFindByID(t *testing.T) {
	repository := NewProductRepository()
	ctx := context.Background()
	product := newTestProduct("product-1", "Keyboard", 120.50, 10)

	if err := repository.Create(ctx, product); err != nil {
		t.Fatalf("expected no error creating product, got %v", err)
	}

	storedProduct, err := repository.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error finding product, got %v", err)
	}

	if storedProduct.ID != product.ID {
		t.Errorf("expected ID %q, got %q", product.ID, storedProduct.ID)
	}

	if storedProduct.Name != product.Name {
		t.Errorf("expected name %q, got %q", product.Name, storedProduct.Name)
	}
}

func TestProductRepositoryRejectsDuplicateID(t *testing.T) {
	repository := NewProductRepository()
	ctx := context.Background()
	product := newTestProduct("product-1", "Keyboard", 120.50, 10)

	if err := repository.Create(ctx, product); err != nil {
		t.Fatalf("expected first creation to succeed, got %v", err)
	}

	err := repository.Create(ctx, product)
	if err == nil {
		t.Fatal("expected duplicate product error")
	}

	var duplicateError *domain.ProductAlreadyExistsError
	if !errors.As(err, &duplicateError) {
		t.Fatalf("expected ProductAlreadyExistsError, got %T", err)
	}
}

func TestProductRepositoryUpdate(t *testing.T) {
	repository := NewProductRepository()
	ctx := context.Background()
	product := newTestProduct("product-1", "Keyboard", 120.50, 10)

	if err := repository.Create(ctx, product); err != nil {
		t.Fatalf("expected no error creating product, got %v", err)
	}

	product.Name = "Mechanical Keyboard"
	product.Price = 150

	if err := repository.Update(ctx, product); err != nil {
		t.Fatalf("expected no error updating product, got %v", err)
	}

	updatedProduct, err := repository.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error finding updated product, got %v", err)
	}

	if updatedProduct.Name != "Mechanical Keyboard" {
		t.Errorf("expected updated name, got %q", updatedProduct.Name)
	}

	if updatedProduct.Price != 150 {
		t.Errorf("expected updated price 150, got %v", updatedProduct.Price)
	}
}

func TestProductRepositoryDelete(t *testing.T) {
	repository := NewProductRepository()
	ctx := context.Background()
	product := newTestProduct("product-1", "Keyboard", 120.50, 10)

	if err := repository.Create(ctx, product); err != nil {
		t.Fatalf("expected no error creating product, got %v", err)
	}

	if err := repository.Delete(ctx, product.ID); err != nil {
		t.Fatalf("expected no error deleting product, got %v", err)
	}

	_, err := repository.FindByID(ctx, product.ID)
	if err == nil {
		t.Fatal("expected product not found error")
	}

	var notFoundError *domain.ProductNotFoundError
	if !errors.As(err, &notFoundError) {
		t.Fatalf("expected ProductNotFoundError, got %T", err)
	}
}

func TestProductRepositoryReturnsCopies(t *testing.T) {
	repository := NewProductRepository()
	ctx := context.Background()
	product := newTestProduct("product-1", "Keyboard", 120.50, 10)

	if err := repository.Create(ctx, product); err != nil {
		t.Fatalf("expected no error creating product, got %v", err)
	}

	product.Name = "Changed outside repository"

	storedProduct, err := repository.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error finding product, got %v", err)
	}

	if storedProduct.Name != "Keyboard" {
		t.Errorf("stored product was modified externally: %q", storedProduct.Name)
	}

	storedProduct.Name = "Changed returned copy"

	storedAgain, err := repository.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error finding product again, got %v", err)
	}

	if storedAgain.Name != "Keyboard" {
		t.Errorf("repository state was modified through returned copy: %q", storedAgain.Name)
	}
}

func newTestProduct(
	id string,
	name string,
	price float64,
	stock int,
) *domain.Product {
	return &domain.Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: time.Date(2026, time.July, 21, 12, 0, 0, 0, time.UTC),
	}
}
