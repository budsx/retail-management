package model

import "time"

type Product struct {
	ProductID   int64       `json:"product_id"`
	ProductName string    `json:"product_name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required"`
	SKU         string    `json:"sku" validate:"required,unique"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}