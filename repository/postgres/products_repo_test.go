package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/budsx/retail-management/model"
	"github.com/stretchr/testify/assert"
)

func Test_ReadProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name    string
		id      int64
		mock    func(sqlmock.Sqlmock)
		want    model.Product
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully get product",
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"product_id", "product_name", "description", "price", "sku", "created_at", "updated_at",
				}).AddRow(1, "Test Product", "Description", 100.0, "SKU123", fixedTime, fixedTime)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, product_name, description, price, sku, created_at, updated_at FROM mst_product WHERE product_id = $1`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: model.Product{
				ProductID:   1,
				ProductName: "Test Product",
				Description: "Description",
				Price:       100.0,
				SKU:         "SKU123",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			wantErr: false,
		},
		{
			name: "Product not found",
			id:   999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, product_name, description, price, sku, created_at, updated_at FROM mst_product WHERE product_id = $1`)).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			want:    model.Product{},
			wantErr: true,
			errMsg:  "product with id 999 not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.ReadProductByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ReadProductsWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	fixedTime := time.Now()

	tests := []struct {
		name    string
		limit   int32
		offset  int32
		mock    func(sqlmock.Sqlmock)
		want    []model.Product
		wantErr bool
	}{
		{
			name:   "Successfully get products",
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"product_id", "product_name", "description", "price", "sku", "created_at", "updated_at",
				}).
					AddRow(1, "Product 1", "Desc 1", 100.0, "SKU1", fixedTime, fixedTime).
					AddRow(2, "Product 2", "Desc 2", 200.0, "SKU2", fixedTime, fixedTime)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, product_name, description, price, sku, created_at, updated_at FROM mst_product ORDER BY product_id LIMIT $1 OFFSET $2`)).
					WithArgs(int32(10), int32(0)).
					WillReturnRows(rows)
			},
			want: []model.Product{
				{
					ProductID:   1,
					ProductName: "Product 1",
					Description: "Desc 1",
					Price:       100.0,
					SKU:         "SKU1",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
				{
					ProductID:   2,
					ProductName: "Product 2",
					Description: "Desc 2",
					Price:       200.0,
					SKU:         "SKU2",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
			},
			wantErr: false,
		},
		{
			name:   "Successfully get products",
			limit:  10,
			offset: 100,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"product_id", "product_name", "description", "price", "sku", "created_at", "updated_at",
				})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_id, product_name, description, price, sku, created_at, updated_at FROM mst_product ORDER BY product_id LIMIT $1 OFFSET $2`)).
					WithArgs(int32(10), int32(100)).
					WillReturnRows(rows)
			},
			want:    []model.Product{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			got, err := rw.ReadProductsWithPagination(context.Background(), tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_UpdateProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		product model.Product
		mock    func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			product: model.Product{
				ProductID:   1,
				ProductName: "Updated Product",
				Description: "Updated Description",
				Price:       150.0,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_product SET product_name = $1, description = $2, price = $3, updated_at = CURRENT_TIMESTAMP WHERE product_id = $4`)).
					WithArgs("Updated Product", "Updated Description", 150.0, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "product not found",
			product: model.Product{
				ProductID:   999,
				ProductName: "Updated Product",
				Description: "Updated Description",
				Price:       150.0,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mst_product SET product_name = $1, description = $2, price = $3, updated_at = CURRENT_TIMESTAMP WHERE product_id = $4`)).
					WithArgs("Updated Product", "Updated Description", 150.0, 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  "product with id 999 not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.UpdateProductByID(context.Background(), tt.product)
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

func Test_WriteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		product model.Product
		mock    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success",
			product: model.Product{
				ProductName: "New Product",
				Description: "New Description",
				Price:       100.0,
				SKU:         "SKU123",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_product (product_name, description, price, sku, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)).
					WithArgs("New Product", "New Description", 100.0, "SKU123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "duplicate SKU error",
			product: model.Product{
				ProductName: "New Product",
				Description: "New Description",
				Price:       100.0,
				SKU:         "SKU123",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_product (product_name, description, price, sku, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)).
					WithArgs("New Product", "New Description", 100.0, "SKU123").
					WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &dbReadWriter{db: db}
			tt.mock(mock)

			err := rw.WriteProduct(context.Background(), tt.product)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
