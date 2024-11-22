package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetTotalStockByLocation(ctx context.Context, locationID int64) ([]model.ProductStock, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetStockTransactionByID - %+v", user))

	totalStock, err := svc.repo.Postgres.GetTotalStockByLocation(ctx, locationID)
	if err != nil {
		svc.logger.Info(err.Error())
		return nil, fmt.Errorf("failed to retrieve total stock for location %d: %w", locationID, err)
	}
	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", totalStock))
	return totalStock, nil
}

func (svc *Service) GetTotalStocks(ctx context.Context) ([]model.ProductStock, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetTotalStocks - %+v", user))

	totalStock, err := svc.repo.Postgres.GetTotalStocks(ctx)
	if err != nil {
		svc.logger.Info(err.Error())
		return nil, fmt.Errorf("failed to retrieve total stock from all locations: %w", err)
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", totalStock))
	return totalStock, nil
}
