package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (svc *Service) AddProduct(ctx context.Context, product model.Product) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Add new product: %+v", product))

	err := svc.repo.Postgres.WriteProduct(ctx, product)
	if err != nil {
		svc.logger.Info(err.Error())
		return err
	}

	svc.logger.Info("[RESPONSE] Product added successfully")
	return nil
}
