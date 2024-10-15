package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/budsx/retail-management/model"
	_ "github.com/lib/pq"
)

type PostgresRepository interface {
	ReadProductByID(context.Context, int64) (model.Product, error)
	ReadProductsWithPagination(context.Context, int32, int32) ([]model.Product, error)
	UpdateProductByID(context.Context, model.Product) error
	WriteProduct(context.Context, model.Product) error

	// User
	RegisterUser(context.Context, model.User) error
	GetUserByUsername(ctx context.Context, username string) (model.User, error)

	// Location & Warehouse
	WriteLocation(ctx context.Context, location model.Location) error
	UpdateLocation(ctx context.Context, location model.Location) error
	ReadLocationByID(ctx context.Context, locationID int64) (model.Location, error)
	DeleteLocationByUserID(ctx context.Context, userID, locationID int64) error
	WriteWarehouse(ctx context.Context, warehouse model.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse model.Warehouse) error
	ReadWarehousesByUserID(ctx context.Context, userID int64) ([]model.Warehouse, error)
	ReadWarehouseByID(ctx context.Context, warehouseID int64) (model.Warehouse, error)

	CreateStockTransaction(context.Context, model.StockTransaction) error
	GetTotalStockByProductAndWarehouse(context.Context, int64, int64) (int64, error)
	GetStockTransactions(context.Context, int64) ([]model.StockTransaction, error)
	GetStockTransactionByID(ctx context.Context, transactionID int64) (model.StockTransaction, error)

	io.Closer
}

type dbReadWriter struct {
	db *sql.DB
}

func NewPostgres(Host, User, Password, DBName string) (PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, 5432, User, Password, DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &dbReadWriter{db}, nil
}

func (rw *dbReadWriter) Close() error {
	if rw.db != nil {
		if err := rw.db.Close(); err != nil {
			return err
		}
		rw.db = nil
	}
	return nil
}
