package model

import "time"

type Warehouse struct {
	WarehouseID   int64     `json:"warehouse_id"`
	WarehouseName string    `json:"warehouse_name"`
	UserID        int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Location struct {
	LocationID   int64     `json:"location_id"`
	LocationName string    `json:"location_name"`
	WarehouseID  int64     `json:"warehouse_id"`
	CreatedAt    time.Time `json:"created_at"`
}
