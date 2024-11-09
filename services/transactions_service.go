package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (s *Service) CreateStockTransaction(ctx context.Context, transaction model.StockTransaction) error {
	s.logger.Info(fmt.Sprintf("[REQUEST] %+v", transaction))
	totalStock, err := s.repo.Postgres.GetTotalStockByProductAndWarehouse(ctx, transaction.ProductID, transaction.WarehouseID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[ERROR] Failed to GetTotalStockByProductAndWarehouse: %s", err.Error()))
		return fmt.Errorf("failed to fetch stock for validation: %w", err)
	}

	if transaction.TransactionType == model.StockIn {
		transaction.Quantity += totalStock
	} else {
		if totalStock < transaction.Quantity {
			s.logger.Error(fmt.Sprintf("Bad request - Total stock %d - Transaction %d", totalStock, transaction.Quantity))
			return fmt.Errorf("stock quantity cannot be negative")
		}
		transaction.Quantity = totalStock - transaction.Quantity
	}

	err = s.repo.Postgres.CreateStockTransaction(ctx, transaction)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[ERROR] Failed to CreateStockTransaction: %s", err.Error()))
		return fmt.Errorf("failed to create stock transaction: %w", err)
	}

	s.logger.Info("[RESPONSE] Create stock successfully")
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
