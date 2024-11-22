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
)

func Test_CreateStockTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		transaction model.StockTransaction
		mockSetup   func(sqlmock.Sqlmock)
		wantErr     bool
	}{
		{
			name: "Successfully create transaction",
			transaction: model.StockTransaction{
				ProductID:       1,
				WarehouseID:     1,
				TransactionType: "IN",
				Quantity:        10,
				CreatedBy:       1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO trx_stock`)).
					WithArgs(1, 1, "IN", 10, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_stock`)).
					WithArgs(10, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed insert",
			transaction: model.StockTransaction{
				ProductID:       1,
				WarehouseID:     1,
				TransactionType: "IN",
				Quantity:        10,
				CreatedBy:       1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO trx_stock`)).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mockSetup(mock)

			err := rw.CreateStockTransaction(context.Background(), tt.transaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStockTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetTotalStockByProductAndWarehouse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		productID   int64
		warehouseID int64
		mockSetup   func(sqlmock.Sqlmock)
		want        int64
		wantErr     bool
	}{
		{
			name:        "Successfully get total stock",
			productID:   1,
			warehouseID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"stock_quantity"}).AddRow(100)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT stock_quantity FROM mst_stock`)).
					WithArgs(int64(1), int64(1)).
					WillReturnRows(rows)
			},
			want:    100,
			wantErr: false,
		},
		{
			name:        "No stock",
			productID:   1,
			warehouseID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT stock_quantity FROM mst_stock`)).
					WithArgs(int64(1), int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mockSetup(mock)

			got, err := rw.GetTotalStockByProductAndWarehouse(context.Background(), tt.productID, tt.warehouseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalStockByProductAndWarehouse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalStockByProductAndWarehouse() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetStockTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name      string
		userID    int64
		mockSetup func(sqlmock.Sqlmock)
		want      []model.StockTransaction
		wantErr   bool
	}{
		{
			name:   "Successfully get transactions",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"transaction_id", "product_id", "warehouse_id",
					"transaction_type", "quantity", "transaction_date", "created_by",
				}).AddRow(1, 1, 1, "IN", 10, fixedTime, 1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by FROM trx_stock`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			want: []model.StockTransaction{{
				TransactionID:   1,
				ProductID:       1,
				WarehouseID:     1,
				TransactionType: "IN",
				Quantity:        10,
				TransactionDate: fixedTime,
				CreatedBy:       1,
			}},
			wantErr: false,
		},
		{
			name:   "no transactions",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by FROM trx_stock`)).
					WithArgs(int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mockSetup(mock)

			got, err := rw.GetStockTransactions(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetTotalStocks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		want      []model.ProductStock
		wantErr   bool
	}{
		{
			name: "Successfully get total stock",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"product_id", "total_stock", "product_name", "sku",
				}).AddRow(1, 100, "Product 1", "SKU001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT m.product_id, SUM(quantity) as total_stock, m.product_name, m.sku FROM trx_stock`)).
					WillReturnRows(rows)
			},
			want: []model.ProductStock{{
				ProductID:   1,
				TotalStock:  100,
				ProductName: "Product 1",
				SKU:         "SKU001",
			}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mockSetup(mock)

			got, err := rw.GetTotalStocks(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalStocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetStockTransactionByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name          string
		transactionID int64
		mock          func(sqlmock.Sqlmock)
		want          model.StockTransaction
		wantErr       bool
	}{
		{
			name:          "success",
			transactionID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"transaction_id",
					"product_id",
					"warehouse_id",
					"transaction_type",
					"quantity",
					"transaction_date",
					"created_by",
				}).AddRow(1, 1, 1, "IN", 10, fixedTime, 1)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by FROM trx_stock WHERE transaction_id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: model.StockTransaction{
				TransactionID:   1,
				ProductID:      1,
				WarehouseID:   1,
				TransactionType: "IN",
				Quantity:      10,
				TransactionDate: fixedTime,
				CreatedBy:     1,
			},
			wantErr: false,
		},
		{
			name:          "transaction not found",
			transactionID: 999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by FROM trx_stock WHERE transaction_id = $1`)).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			want:    model.StockTransaction{},
			wantErr: true,
		},
		{
			name:          "database error",
			transactionID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by FROM trx_stock WHERE transaction_id = $1`)).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    model.StockTransaction{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.GetStockTransactionByID(context.Background(), tt.transactionID)
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

func Test_GetTotalStockByLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		locationID int64
		mock       func(sqlmock.Sqlmock)
		want       []model.ProductStock
		wantErr    bool
	}{
		{
			name:       "success with multiple products",
			locationID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "total_stock"}).
					AddRow(1, 100).
					AddRow(2, 200)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, SUM(quantity) as total_stock FROM trx_stock WHERE warehouse_id = $1 GROUP BY product_id`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: []model.ProductStock{
				{
					ProductID:   1,
					TotalStock: 100,
				},
				{
					ProductID:   2,
					TotalStock: 200,
				},
			},
			wantErr: false,
		},
		{
			name:       "success with no stock",
			locationID: 2,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "total_stock"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, SUM(quantity) as total_stock FROM trx_stock WHERE warehouse_id = $1 GROUP BY product_id`)).
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    []model.ProductStock{},
			wantErr: false,
		},
		{
			name:       "database error",
			locationID: 3,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, SUM(quantity) as total_stock FROM trx_stock WHERE warehouse_id = $1 GROUP BY product_id`)).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "scan error",
			locationID: 4,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "total_stock"}).
					AddRow("invalid", 100) // This will cause a scan error
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, SUM(quantity) as total_stock FROM trx_stock WHERE warehouse_id = $1 GROUP BY product_id`)).
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

			got, err := rw.GetTotalStockByLocation(context.Background(), tt.locationID)
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
