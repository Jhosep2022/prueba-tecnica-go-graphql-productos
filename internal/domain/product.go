package domain

import "time"

// Product represents a product available in the store.
type Product struct {
	ID        string
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}
