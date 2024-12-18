package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/model"
	"github.com/gorilla/mux"
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

func (c *Controller) GetWarehousesByUserID(w http.ResponseWriter, r *http.Request) {
	warehouses, err := c.service.GetWarehouseByUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, warehouses)
}

func (c *Controller) EditWarehouseByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Missing warehouse ID")
		return
	}

	warehouseID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	var warehouse model.Warehouse
	err = json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	warehouse.WarehouseID = warehouseID
	err = c.service.EditWarehouseByUserID(r.Context(), warehouse)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Warehouse updated successfully")
}
