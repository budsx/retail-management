package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetStockTransactions(ctx context.Context) ([]model.StockTransaction, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetStockTransactions - %+v", user))

	if user.UserID == 0 {
		svc.logger.Info("Invalid User")
		return []model.StockTransaction{}, fmt.Errorf("Unathorized")
	}

	transactions, err := svc.repo.Postgres.GetStockTransactions(ctx, user.UserID)
	if err != nil {
		svc.logger.Info(err.Error())
		return []model.StockTransaction{}, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", transactions))
	return transactions, nil
}
