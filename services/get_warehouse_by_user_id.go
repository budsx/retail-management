package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetWarehouseByUserID(ctx context.Context) ([]model.Warehouse, error) {
	svc.logger.Info("[REQUEST] Get warehouse by user ID")

	userID := ctx.Value(middleware.ContextKeyUserID).(int64)

	warehouses, err := svc.repo.Postgres.ReadWarehousesByUserID(ctx, userID)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to get warehouses: %s", err.Error()))
		return nil, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] Warehouses retrieved successfully: %+v", warehouses))
	return warehouses, nil
}
