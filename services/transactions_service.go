package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) CreateStockTransaction(ctx context.Context, transaction model.StockTransaction) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] %+v", transaction))
	totalStock, err := svc.repo.Postgres.GetTotalStockByProductAndWarehouse(ctx, transaction.ProductID, transaction.WarehouseID)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to GetTotalStockByProductAndWarehouse: %s", err.Error()))
		return fmt.Errorf("failed to fetch stock for validation: %w", err)
	}

	if transaction.TransactionType == model.StockIn {
		transaction.Quantity += totalStock
	} else {
		if totalStock < transaction.Quantity {
			svc.logger.Error(fmt.Sprintf("Bad request - Total stock %d - Transaction %d", totalStock, transaction.Quantity))
			return fmt.Errorf("stock quantity cannot be negative")
		}
		transaction.Quantity = totalStock - transaction.Quantity
	}

	err = svc.repo.Postgres.CreateStockTransaction(ctx, transaction)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to CreateStockTransaction: %s", err.Error()))
		return fmt.Errorf("failed to create stock transaction: %w", err)
	}

	svc.logger.Info("[RESPONSE] Create stock successfully")
	return nil
}

func (svc *Service) GetStockTransactionByID(ctx context.Context, transactionID int64) (model.StockTransaction, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetStockTransactionByID - %+v", user))

	transaction, err := svc.repo.Postgres.GetStockTransactionByID(ctx, transactionID)
	if err != nil {
		svc.logger.Info(err.Error())
		return model.StockTransaction{}, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", transaction))
	return transaction, nil
}
