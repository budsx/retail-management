package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

//  Stock All Location
func (svc *Service) GetTotalStocks(ctx context.Context) ([]model.ProductStock, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetTotalStocks - %+v", user))

	totalStock, err := svc.repo.Postgres.GetTotalStocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total stock from all locations: %w", err)
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", totalStock))
	return totalStock, nil
}
