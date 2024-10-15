package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) AddLocation(ctx context.Context, location model.Location) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Add new location: %+v", location))

	userID := ctx.Value(middleware.ContextKeyUserID).(int64)

	warehouse, err := svc.repo.Postgres.ReadWarehouseByID(ctx, location.WarehouseID)
	if err != nil || warehouse.UserID != userID {
		svc.logger.Error("[ERROR] Unauthorized or warehouse not found")
		return fmt.Errorf("unauthorized or warehouse not found")
	}

	err = svc.repo.Postgres.WriteLocation(ctx, location)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to add location: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] Location added successfully")
	return nil
}