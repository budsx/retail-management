package services

import (
	"context"

	"github.com/budsx/retail-management/model"
	"github.com/budsx/retail-management/repository"
	"github.com/budsx/retail-management/utils"
)

type RetailManagementService interface {
	GetProductByID(context.Context, int64) (model.Product, error)
	GetProducts(context.Context, model.Pagination) ([]model.Product, error)
	AddProduct(context.Context, model.Product) (error)
	EditProduct(context.Context, model.Product) error
}

type Service struct {
	repo repository.Repository
	logger utils.Interface
}

func NewRetailManagementService(repo repository.Repository, logger utils.Interface) RetailManagementService {
	return &Service{repo: repo, logger: logger}
}
