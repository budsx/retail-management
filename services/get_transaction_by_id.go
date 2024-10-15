package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

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
