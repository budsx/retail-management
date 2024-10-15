package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetWarehouseByUserID(ctx context.Context) ([]model.Warehouse, error) {
	user := middleware.GetUserInfoByContext(ctx)

	svc.logger.Info(fmt.Sprintf("[REQUEST] Get warehouse by user ID - %+v", user))

	warehouses, err := svc.repo.Postgres.ReadWarehousesByUserID(ctx, user.UserID)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to get warehouses: %s", err.Error()))
		return nil, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] Warehouses retrieved successfully: %+v", warehouses))
	return warehouses, nil
}
