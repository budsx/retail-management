package postgres

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (rw *dbReadWriter) ReadWarehouseByID(ctx context.Context, warehouseID int64) (model.Warehouse, error) {
	selectWarehouseByID := `SELECT warehouse_id, warehouse_name, user_id, created_at 
							FROM mst_warehouse WHERE warehouse_id = $1`

	var warehouse model.Warehouse
	err := rw.db.QueryRowContext(ctx, selectWarehouseByID, warehouseID).Scan(&warehouse.WarehouseID, &warehouse.WarehouseName, &warehouse.UserID, &warehouse.CreatedAt)
	if err != nil {
		return model.Warehouse{}, err
	}

	return warehouse, nil
}

func (rw *dbReadWriter) WriteWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	insertWarehouse := `INSERT INTO mst_warehouse (warehouse_name, user_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`

	_, err := rw.db.ExecContext(ctx, insertWarehouse, warehouse.WarehouseName, warehouse.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) UpdateWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	updateWarehouse := `UPDATE mst_warehouse SET warehouse_name = $1 WHERE warehouse_id = $2`

	_, err := rw.db.ExecContext(ctx, updateWarehouse, warehouse.WarehouseName, warehouse.WarehouseID)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) WriteLocation(ctx context.Context, location model.Location) error {
	insertLocation := `INSERT INTO mst_location (location_name, warehouse_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`

	_, err := rw.db.ExecContext(ctx, insertLocation, location.LocationName, location.WarehouseID)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) ReadLocationByID(ctx context.Context, locationID int64) (model.Location, error) {
	selectLocationByID := `SELECT location_id, location_name, warehouse_id, created_at 
						   FROM mst_location WHERE location_id = $1`

	var location model.Location
	err := rw.db.QueryRowContext(ctx, selectLocationByID, locationID).Scan(&location.LocationID, &location.LocationName, &location.WarehouseID, &location.CreatedAt)
	if err != nil {
		return model.Location{}, err
	}

	return location, err
}

func (rw *dbReadWriter) UpdateLocation(ctx context.Context, location model.Location) error {
	updateLocation := `UPDATE mst_location SET location_name = $1 WHERE location_id = $2`

	_, err := rw.db.ExecContext(ctx, updateLocation, location.LocationName, location.LocationID)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) DeleteLocationByUserID(ctx context.Context, userID, locationID int64) error {
	selectLocation := `
		SELECT l.location_id
		FROM mst_location l
		INNER JOIN mst_warehouse w ON l.warehouse_id = w.warehouse_id
		WHERE l.location_id = $1 AND w.user_id = $2`

	var locationIDCheck int64
	err := rw.db.QueryRowContext(ctx, selectLocation, locationID, userID).Scan(&locationIDCheck)
	if err != nil {
		return fmt.Errorf("unauthorized or location not found")
	}

	deleteLocation := `DELETE FROM mst_location WHERE location_id = $1`
	_, err = rw.db.ExecContext(ctx, deleteLocation, locationID)
	if err != nil {
		return fmt.Errorf("failed to delete location: %v", err)
	}

	return nil
}

func (rw *dbReadWriter) ReadWarehousesByUserID(ctx context.Context, userID int64) ([]model.Warehouse, error) {
	warehouses := make([]model.Warehouse, 0)

	selectWarehouseByUser := `SELECT warehouse_id, warehouse_name, user_id, created_at 
	          FROM mst_warehouse 
	          WHERE user_id = $1`

	rows, err := rw.db.QueryContext(ctx, selectWarehouseByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse model.Warehouse
		err := rows.Scan(&warehouse.WarehouseID, &warehouse.WarehouseName, &warehouse.UserID, &warehouse.CreatedAt)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return warehouses, nil
}
