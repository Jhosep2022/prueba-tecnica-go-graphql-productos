package domain

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

type ProductNotFoundError struct {
	ID string
}

func (e *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product with id %q was not found", e.ID)
}
