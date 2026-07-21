package graphql

import (
	"context"
	"errors"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/Jhosep2022/prueba-tecnica-go-graphql-productos/internal/domain"
)

const (
	codeValidationError      = "VALIDATION_ERROR"
	codeProductNotFound      = "PRODUCT_NOT_FOUND"
	codeProductAlreadyExists = "PRODUCT_ALREADY_EXISTS"
	codeGraphQLError         = "GRAPHQL_ERROR"
	codeInternalError        = "INTERNAL_ERROR"
)

// ErrorPresenter maps application errors to safe GraphQL errors.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	presentedError := gqlgraphql.DefaultErrorPresenter(ctx, err)

	var validationError *domain.ValidationError
	if errors.As(err, &validationError) {
		presentedError.Message = validationError.Error()
		presentedError.Extensions = map[string]any{
			"code":  codeValidationError,
			"field": validationError.Field,
		}

		return presentedError
	}

	var notFoundError *domain.ProductNotFoundError
	if errors.As(err, &notFoundError) {
		presentedError.Message = notFoundError.Error()
		presentedError.Extensions = map[string]any{
			"code": codeProductNotFound,
		}

		return presentedError
	}

	var alreadyExistsError *domain.ProductAlreadyExistsError
	if errors.As(err, &alreadyExistsError) {
		presentedError.Message = alreadyExistsError.Error()
		presentedError.Extensions = map[string]any{
			"code": codeProductAlreadyExists,
		}

		return presentedError
	}

	var graphQLError *gqlerror.Error
	if errors.As(err, &graphQLError) {
		if presentedError.Extensions == nil {
			presentedError.Extensions = make(map[string]any)
		}

		if _, exists := presentedError.Extensions["code"]; !exists {
			presentedError.Extensions["code"] = codeGraphQLError
		}

		return presentedError
	}

	presentedError.Message = "internal server error"
	presentedError.Extensions = map[string]any{
		"code": codeInternalError,
	}

	return presentedError
}
