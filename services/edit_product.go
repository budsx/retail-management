package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (svc *Service) EditProduct(ctx context.Context, updatedProduct model.Product) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Update product: %+v", updatedProduct))

	err := svc.repo.Postgres.UpdateProductByID(ctx, updatedProduct)
	if err != nil {
		svc.logger.Info(err.Error())
		return err
	}

	svc.logger.Info("[RESPONSE] Product updated successfully")
	return nil
}
