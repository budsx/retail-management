package controller

import (
	"net/http"

	"github.com/budsx/retail-management/services"
)

type RetailManagementController interface {
	// Health Check
	Health(w http.ResponseWriter, r *http.Request)
	// Users

	// Product
	GetProductByID(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
	EditProduct(w http.ResponseWriter, r *http.Request)
	AddProduct(w http.ResponseWriter, r *http.Request)

	// Transaction
}

type Controller struct {
	service services.RetailManagementService
}

func NewRetailManagementController(service services.RetailManagementService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
