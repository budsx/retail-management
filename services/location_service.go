package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (svc *Service) AddLocation(ctx context.Context, location model.Location) error {
	user := middleware.GetUserInfoByContext(ctx)

	svc.logger.Info(fmt.Sprintf("[REQUEST] Add new location: %+v - %+v", location, user))

	warehouse, err := svc.repo.Postgres.ReadWarehouseByID(ctx, location.WarehouseID)
	if err != nil || warehouse.UserID != user.UserID {
		svc.logger.Error("[ERROR] Unauthorized or warehouse not found")
		return fmt.Errorf("unauthorized or warehouse not found")
	}

	err = svc.repo.Postgres.WriteLocation(ctx, location)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to add location: %s", err.Error()))
		return fmt.Errorf("failed to add location: %w", err)
	}

	svc.logger.Info("[RESPONSE] Location added successfully")
	return nil
}

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
		return fmt.Errorf("failed to delete location: %w", err)
	}

	svc.logger.Info("[RESPONSE] Location deleted successfully")
	return nil
}

func (svc *Service) EditLocationByUserID(ctx context.Context, location model.Location) error {
	userID := ctx.Value(middleware.ContextKeyUserID).(int64)
	svc.logger.Info(fmt.Sprintf("[REQUEST]  Edit location: %+v - %d", location, userID))

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
