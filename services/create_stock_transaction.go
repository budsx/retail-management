package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (s *Service) CreateStockTransaction(ctx context.Context, transaction model.StockTransaction) error {
	totalStock, err := s.repo.Postgres.GetTotalStockByProductAndWarehouse(ctx, transaction.ProductID, transaction.WarehouseID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock for validation: %w", err)
	}

	if transaction.TransactionType == model.StockOut && totalStock < transaction.Quantity {
		return fmt.Errorf("insufficient stock")
	}

	err = s.repo.Postgres.CreateStockTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create stock transaction: %w", err)
	}

	return nil
}
