package model

import "time"

type TransactionType string

const (
	StockIn  = TransactionType("IN")
	StockOut = TransactionType("OUT")
)

type StockTransaction struct {
	TransactionID   int64           `json:"transaction_id"`
	ProductID       int64           `json:"product_id"`
	WarehouseID     int64           `json:"warehouse_id"`
	TransactionType TransactionType `json:"transaction_type"`
	Quantity        int64           `json:"quantity"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedBy       int64           `json:"created_by"`
}
