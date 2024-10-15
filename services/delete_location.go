package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
)

func (svc *Service) DeleteLocationByUserID(ctx context.Context, locationID int64) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Delete location ID: %d", locationID))

	userID := ctx.Value(middleware.ContextKeyUserID).(int64)

	dbLocation, err := svc.repo.Postgres.ReadLocationByID(ctx, locationID)
	if err != nil {
		svc.logger.Error("[ERROR] Location not found")
		return fmt.Errorf("location not found")
	}

	warehouse, err := svc.repo.Postgres.ReadWarehouseByID(ctx, dbLocation.WarehouseID)
	if err != nil || warehouse.UserID != userID {
		svc.logger.Error("[ERROR] Unauthorized or location not found")
		return fmt.Errorf("unauthorized or location not found")
	}

	err = svc.repo.Postgres.DeleteLocationByUserID(ctx, userID, locationID)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to delete location: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] Location deleted successfully")
	return nil
}
