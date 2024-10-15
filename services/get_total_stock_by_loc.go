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
		return nil, fmt.Errorf("failed to retrieve total stock for location %d: %w", locationID, err)
	}
	return totalStock, nil
}
