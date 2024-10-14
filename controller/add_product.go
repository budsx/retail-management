package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/model"
)

func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.service.AddProduct(r.Context(), product)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Product added successfully")
}
