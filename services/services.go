package services

import "github.com/budsx/retail-management/repository"

type RetailManagementService interface {
}

type Service struct {
	repo repository.Repository
}

func NewRetailManagementService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}
