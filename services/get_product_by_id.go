package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetProductByID(ctx context.Context, req int64) (model.Product, error) {
	svc.logger.Info(fmt.Sprintf("[REQUEST] %d", req))

	product, err := svc.repo.Postgres.ReadProductByID(ctx, req)
	if err != nil {
		svc.logger.Info(err.Error())
		return model.Product{}, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", product))
	return product, nil
}