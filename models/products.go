package models

import (
	"time"
)

// Product - Product for the app.
type Product struct {
	ID     int64  `json:"id,omitempty"`
	UserID int64  `json:"userID,omitempty"`
	Offer  *Offer `json:"offer,omitempty"`

	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Price       float64    `json:"price,omitempty"`
	Categories  Categories `json:"categories,omitempty"`
	Images      Files      `json:"images,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Products - Alias for a product array.
type Products []Product

// Validate - Validates a product.
func (p Product) Validate() (valid bool) {
	valid = true

	if p.Name == "" {
		valid = false
	}

	if p.UserID == 0 || p.Price == 0 {
		valid = false
	}
	return
}

// GetIdentifier gets the first unique identifier it finds in order of importance.
func (p Product) GetIdentifier() (identifier interface{}) {
	if p.ID != 0 {
		identifier = p.ID
	}

	return
}
