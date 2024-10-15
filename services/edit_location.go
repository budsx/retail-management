package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) EditLocationByUserID(ctx context.Context, location model.Location) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Edit location: %+v", location))

	userID := ctx.Value(middleware.ContextKeyUserID).(int64)

	dbLocation, err := svc.repo.Postgres.ReadLocationByID(ctx, location.LocationID)
	if err != nil {
		svc.logger.Error("[ERROR] Location not found")
		return fmt.Errorf("location not found")
	}

	warehouse, err := svc.repo.Postgres.ReadWarehouseByID(ctx, dbLocation.WarehouseID)
	if err != nil || warehouse.UserID != userID {
		svc.logger.Error("[ERROR] Unauthorized or location not found")
		return fmt.Errorf("unauthorized or location not found")
	}

	err = svc.repo.Postgres.UpdateLocation(ctx, location)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to update location: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] Location updated successfully")
	return nil
}
