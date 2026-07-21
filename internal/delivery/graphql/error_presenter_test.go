package graphql

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

func TestErrorPresenterMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedCode  string
		expectedField string
	}{
		{
			name: "validation error",
			err: &domain.ValidationError{
				Field:   "price",
				Message: "must be greater than zero",
			},
			expectedCode:  codeValidationError,
			expectedField: "price",
		},
		{
			name: "wrapped product not found error",
			err: fmt.Errorf(
				"get product: %w",
				&domain.ProductNotFoundError{ID: "product-1"},
			),
			expectedCode: codeProductNotFound,
		},
		{
			name: "product already exists error",
			err: &domain.ProductAlreadyExistsError{
				ID: "product-1",
			},
			expectedCode: codeProductAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			presentedError := ErrorPresenter(context.Background(), test.err)

			if presentedError.Extensions["code"] != test.expectedCode {
				t.Errorf(
					"expected code %q, got %v",
					test.expectedCode,
					presentedError.Extensions["code"],
				)
			}

			if test.expectedField != "" && presentedError.Extensions["field"] != test.expectedField {
				t.Errorf(
					"expected field %q, got %v",
					test.expectedField,
					presentedError.Extensions["field"],
				)
			}
		})
	}
}

func TestErrorPresenterHidesInternalErrorDetails(t *testing.T) {
	presentedError := ErrorPresenter(
		context.Background(),
		errors.New("database password leaked"),
	)

	if presentedError.Message != "internal server error" {
		t.Errorf("expected safe internal message, got %q", presentedError.Message)
	}

	if presentedError.Extensions["code"] != codeInternalError {
		t.Errorf(
			"expected code %q, got %v",
			codeInternalError,
			presentedError.Extensions["code"],
		)
	}
}
