package services

import (
	"context"

	"github.com/budsx/retail-management/model"
)

func (s *Service) GetStockTransactionByID(ctx context.Context, transactionID int64) (model.StockTransaction, error) {
	transaction, err := s.repo.Postgres.GetStockTransactionByID(ctx, transactionID)
	if err != nil {
		return model.StockTransaction{}, err
	}
	return transaction, nil
}
