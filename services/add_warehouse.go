package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) AddWarehouseByUserID(ctx context.Context, warehouse model.Warehouse) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Add new warehouse: %+v", warehouse))
	userID := ctx.Value(middleware.ContextKeyUserID).(int64)
	warehouse.UserID = userID

	err := svc.repo.Postgres.WriteWarehouse(ctx, warehouse)
	if err != nil {
		return err
	}

	svc.logger.Info("[RESPONSE] Product added successfully")
	return nil
}
