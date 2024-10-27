package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/model"
)

func (c *Controller) AddWarehouseByUserID(w http.ResponseWriter, r *http.Request) {
	var warehouse model.Warehouse
	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.service.AddWarehouseByUserID(r.Context(), warehouse)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Warehouse added successfully")
}
