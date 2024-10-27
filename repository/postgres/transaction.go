package postgres

import (
	"context"

	"github.com/budsx/retail-management/model"
)

func (rw *dbReadWriter) CreateStockTransaction(ctx context.Context, transaction model.StockTransaction) error {
	stockAdjustment := `INSERT INTO trx_stock (product_id, warehouse_id, transaction_type, quantity, created_by) 
              VALUES ($1, $2, $3, $4, $5)`

	updateStock := `UPDATE mst_stock SET stock_quantity = $1 WHERE product_id = $2 AND warehouse_id = $3`

	tx, err := rw.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stockAdjustment, transaction.ProductID, transaction.WarehouseID, transaction.TransactionType, transaction.Quantity, transaction.CreatedBy)
	if err != nil {
		return err
	}

	_, err = tx.Exec(updateStock, transaction.Quantity, transaction.ProductID, transaction.WarehouseID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (rw *dbReadWriter) GetTotalStockByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (int64, error) {
	selectTotalStock := `SELECT stock_quantity FROM mst_stock WHERE product_id = $1 AND warehouse_id = $2`

	var totalStock int64
	err := rw.db.QueryRowContext(ctx, selectTotalStock, productID, warehouseID).Scan(&totalStock)
	if err != nil {
		return 0, err
	}

	return totalStock, nil
}

func (rw *dbReadWriter) GetStockTransactions(ctx context.Context, userID int64) ([]model.StockTransaction, error) {
	selectAllTransaction := `SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by 
	          FROM trx_stock WHERE created_by = $1`

	rows, err := rw.db.QueryContext(ctx, selectAllTransaction, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.StockTransaction
	for rows.Next() {
		var transaction model.StockTransaction
		err := rows.Scan(&transaction.TransactionID, &transaction.ProductID, &transaction.WarehouseID, &transaction.TransactionType, &transaction.Quantity, &transaction.TransactionDate, &transaction.CreatedBy)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (rw *dbReadWriter) GetStockTransactionByID(ctx context.Context, transactionID int64) (model.StockTransaction, error) {
	var transaction model.StockTransaction

	query := `SELECT transaction_id, product_id, warehouse_id, transaction_type, quantity, transaction_date, created_by 
	          FROM trx_stock WHERE transaction_id = $1`
	err := rw.db.QueryRowContext(ctx, query, transactionID).Scan(&transaction.TransactionID, &transaction.ProductID, &transaction.WarehouseID, &transaction.TransactionType, &transaction.Quantity, &transaction.TransactionDate, &transaction.CreatedBy)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (rw *dbReadWriter) GetTotalStocks(ctx context.Context) ([]model.ProductStock, error) {
	query := `SELECT product_id, SUM(quantity) as total_stock
	          FROM trx_stock
	          GROUP BY product_id`

	rows, err := rw.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalStock []model.ProductStock
	for rows.Next() {
		var productStock model.ProductStock
		err := rows.Scan(&productStock.ProductID, &productStock.TotalStock)
		if err != nil {
			return nil, err
		}
		totalStock = append(totalStock, productStock)
	}
	return totalStock, nil
}

func (rw *dbReadWriter) GetTotalStockByLocation(ctx context.Context, locationID int64) ([]model.ProductStock, error) {
	query := `SELECT product_id, SUM(quantity) as total_stock
	          FROM trx_stock
	          WHERE warehouse_id = $1
	          GROUP BY product_id`

	rows, err := rw.db.QueryContext(ctx, query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalStock []model.ProductStock
	for rows.Next() {
		var productStock model.ProductStock
		err := rows.Scan(&productStock.ProductID, &productStock.TotalStock)
		if err != nil {
			return nil, err
		}
		totalStock = append(totalStock, productStock)
	}
	return totalStock, nil
}
