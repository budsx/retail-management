package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) EditWarehouseByUserID(ctx context.Context, warehouse model.Warehouse) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Edit warehouse: %+v", warehouse))

	userID := ctx.Value(middleware.ContextKeyUserID).(int64)

	dbWarehouse, err := svc.repo.Postgres.ReadWarehouseByID(ctx, warehouse.WarehouseID)
	if err != nil || dbWarehouse.UserID != userID {
		svc.logger.Error("[ERROR] Unauthorized or warehouse not found")
		return fmt.Errorf("unauthorized or warehouse not found")
	}

	err = svc.repo.Postgres.UpdateWarehouse(ctx, warehouse)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to update warehouse: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] Warehouse updated successfully")
	return nil
}
