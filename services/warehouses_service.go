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

func (svc *Service) GetStockTransactions(ctx context.Context) ([]model.StockTransaction, error) {
	user := middleware.GetUserInfoByContext(ctx)
	svc.logger.Info(fmt.Sprintf("[REQUEST] GetStockTransactions - %+v", user))

	if user.UserID == 0 {
		svc.logger.Info("Invalid User")
		return []model.StockTransaction{}, fmt.Errorf("Unathorized")
	}

	transactions, err := svc.repo.Postgres.GetStockTransactions(ctx, user.UserID)
	if err != nil {
		svc.logger.Info(err.Error())
		return []model.StockTransaction{}, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", transactions))
	return transactions, nil
}
