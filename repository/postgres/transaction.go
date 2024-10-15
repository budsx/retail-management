package postgres

import (
	"context"

	"github.com/budsx/retail-management/model"
)

func (rw *dbReadWriter) CreateStockTransaction(ctx context.Context, transaction model.StockTransaction) error {
	stockAdjustment := `INSERT INTO trx_stock (product_id, warehouse_id, transaction_type, quantity, created_by) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := rw.db.ExecContext(ctx, stockAdjustment, transaction.ProductID, transaction.WarehouseID, transaction.TransactionType, transaction.Quantity, transaction.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) GetTotalStockByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (int64, error) {
	selectTotalStock := `SELECT SUM(quantity) FROM trx_stock WHERE product_id = $1 AND warehouse_id = $2`

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