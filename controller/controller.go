package controller

import (
	"net/http"

	"github.com/budsx/retail-management/services"
)

type Controller struct {
	service services.RetailManagementService
}

func NewRetailManagementController(service services.RetailManagementService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
