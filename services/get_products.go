package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (svc *Service) GetProducts(ctx context.Context, pagination model.Pagination) ([]model.Product, error) {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Get products with pagination: %+v", pagination))

	offset := (pagination.Page - 1) * pagination.Limit

	products, err := svc.repo.Postgres.ReadProductsWithPagination(ctx, pagination.Limit, offset)
	if err != nil {
		svc.logger.Info(err.Error())
		return nil, err
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] %+v", products))
	return products, nil
}
