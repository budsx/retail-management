package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/budsx/retail-management/model"
	"github.com/stretchr/testify/assert"
	"github.com/lib/pq"
)

func Test_ReadWarehouseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name        string
		warehouseID int64
		mock        func(sqlmock.Sqlmock)
		want        model.Warehouse
		wantErr     bool
	}{
		{
			name:        "success",
			warehouseID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"warehouse_id", "warehouse_name", "user_id", "created_at"}).
					AddRow(1, "Test Warehouse", 1, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE warehouse_id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: model.Warehouse{
				WarehouseID:   1,
				WarehouseName: "Test Warehouse",
				UserID:        1,
				CreatedAt:     fixedTime,
			},
			wantErr: false,
		},
		{
			name:        "not found",
			warehouseID: 999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE warehouse_id = $1`)).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			want:    model.Warehouse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.ReadWarehouseByID(context.Background(), tt.warehouseID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_WriteWarehouse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name      string
		warehouse model.Warehouse
		mock      func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			warehouse: model.Warehouse{
				WarehouseName: "New Warehouse",
				UserID:        1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_warehouse (warehouse_name, user_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("New Warehouse", int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.WriteWarehouse(context.Background(), tt.warehouse)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_ReadLocationByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name       string
		locationID int64
		mock       func(sqlmock.Sqlmock)
		want       model.Location
		wantErr    bool
	}{
		{
			name:       "success",
			locationID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"location_id", "location_name", "warehouse_id", "created_at"}).
					AddRow(1, "Test Location", 1, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT location_id, location_name, warehouse_id, created_at FROM mst_location WHERE location_id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: model.Location{
				LocationID:   1,
				LocationName: "Test Location",
				WarehouseID:  1,
				CreatedAt:    fixedTime,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.ReadLocationByID(context.Background(), tt.locationID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_DeleteLocationByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		userID     int64
		locationID int64
		mock       func(sqlmock.Sqlmock)
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "success",
			userID:     1,
			locationID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.location_id FROM mst_location l INNER JOIN mst_warehouse w`)).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"location_id"}).AddRow(1))

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM mst_location WHERE location_id = $1`)).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:       "unauthorized",
			userID:     2,
			locationID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.location_id FROM mst_location l INNER JOIN mst_warehouse w`)).
					WithArgs(1, 2).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  "unauthorized or location not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.DeleteLocationByUserID(context.Background(), tt.userID, tt.locationID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_ReadWarehousesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name    string
		userID  int64
		mock    func(sqlmock.Sqlmock)
		want    []model.Warehouse
		wantErr bool
	}{
		{
			name:   "success with multiple warehouses",
			userID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"warehouse_id", "warehouse_name", "user_id", "created_at"}).
					AddRow(1, "Warehouse 1", 1, fixedTime).
					AddRow(2, "Warehouse 2", 1, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE user_id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: []model.Warehouse{
				{
					WarehouseID:   1,
					WarehouseName: "Warehouse 1",
					UserID:        1,
					CreatedAt:     fixedTime,
				},
				{
					WarehouseID:   2,
					WarehouseName: "Warehouse 2",
					UserID:        1,
					CreatedAt:     fixedTime,
				},
			},
			wantErr: false,
		},
		{
			name:   "success with no warehouses",
			userID: 2,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"warehouse_id", "warehouse_name", "user_id", "created_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE user_id = $1`)).
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    []model.Warehouse{},
			wantErr: false,
		},
		{
			name:   "database error",
			userID: 3,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE user_id = $1`)).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "error during row scan",
			userID: 4,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"warehouse_id", "warehouse_name", "user_id", "created_at"}).
					AddRow("invalid", "Warehouse 1", 1, fixedTime) // warehouse_id as string will cause scan error
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT warehouse_id, warehouse_name, user_id, created_at FROM mst_warehouse WHERE user_id = $1`)).
					WithArgs(4).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.ReadWarehousesByUserID(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_UpdateLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		location model.Location
		mock     func(sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name: "success",
			location: model.Location{
				LocationID:   1,
				LocationName: "Updated Location",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_location SET location_name = $1 WHERE location_id = $2`)).
					WithArgs("Updated Location", int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "location not found",
			location: model.Location{
				LocationID:   999,
				LocationName: "Non-existent Location",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_location SET location_name = $1 WHERE location_id = $2`)).
					WithArgs("Non-existent Location", int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name: "database error",
			location: model.Location{
				LocationID:   1,
				LocationName: "Error Location",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_location SET location_name = $1 WHERE location_id = $2`)).
					WithArgs("Error Location", int64(1)).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.UpdateLocation(context.Background(), tt.location)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_UpdateWarehouse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name      string
		warehouse model.Warehouse
		mock      func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			warehouse: model.Warehouse{
				WarehouseID:   1,
				WarehouseName: "Updated Warehouse",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_warehouse SET warehouse_name = $1 WHERE warehouse_id = $2`)).
					WithArgs("Updated Warehouse", int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "warehouse not found",
			warehouse: model.Warehouse{
				WarehouseID:   999,
				WarehouseName: "Non-existent Warehouse",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_warehouse SET warehouse_name = $1 WHERE warehouse_id = $2`)).
					WithArgs("Non-existent Warehouse", int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name: "database error",
			warehouse: model.Warehouse{
				WarehouseID:   1,
				WarehouseName: "Error Warehouse",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_warehouse SET warehouse_name = $1 WHERE warehouse_id = $2`)).
					WithArgs("Error Warehouse", int64(1)).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.UpdateWarehouse(context.Background(), tt.warehouse)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_WriteLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		location model.Location
		mock     func(sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name: "success",
			location: model.Location{
				LocationName: "New Location",
				WarehouseID: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_location (location_name, warehouse_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("New Location", int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "foreign key violation",
			location: model.Location{
				LocationName: "Invalid Location",
				WarehouseID: 999, // Non-existent warehouse
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_location (location_name, warehouse_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("Invalid Location", int64(999)).
					WillReturnError(&pq.Error{Code: "23503"}) // Foreign key violation
			},
			wantErr: true,
		},
		{
			name: "database error",
			location: model.Location{
				LocationName: "Error Location",
				WarehouseID: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_location (location_name, warehouse_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("Error Location", int64(1)).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.WriteLocation(context.Background(), tt.location)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
