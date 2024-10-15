package postgres

import (
	"context"
	"io"

	"github.com/budsx/retail-management/model"
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
	GetTotalStocks(ctx context.Context) ([]model.ProductStock, error)
	GetTotalStockByLocation(context.Context, int64) ([]model.ProductStock, error)

	io.Closer
}
