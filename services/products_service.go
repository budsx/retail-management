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
