package controller

import (
	"net/http"

	"github.com/budsx/retail-management/services"
)

type RetailManagementController interface {
	// Product
	GetProductByID(w http.ResponseWriter, r *http.Request)

	// Stock
	GetStock(w http.ResponseWriter, r *http.Request)
}

type Controller struct {
	service services.RetailManagementService
}

func NewRetailManagementController(service services.RetailManagementService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetProductByID(w http.ResponseWriter, r *http.Request) {
	return
}

func (c *Controller) GetStock(w http.ResponseWriter, r *http.Request) {
	return
}