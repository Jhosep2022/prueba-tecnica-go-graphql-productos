package domain

import "fmt"

// ValidationError is returned when a validation check fails for a specific field.
type ValidationError struct {
	Field   string
	Message string
}

// Error returns a string representation of the ValidationError.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// ProductNotFoundError is returned when a product with a specific ID is not found in the repository.
type ProductNotFoundError struct {
	ID string
}

// Error returns a string representation of the ProductNotFoundError.
func (e *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product with id %q was not found", e.ID)
}

// ProductAlreadyExistsError is returned when trying to create a product with an ID that already exists in the repository.
type ProductAlreadyExistsError struct {
	ID string
}

// Error returns a string representation of the ProductAlreadyExistsError.
func (e *ProductAlreadyExistsError) Error() string {
	return fmt.Sprintf("product with id %q already exists", e.ID)
}
