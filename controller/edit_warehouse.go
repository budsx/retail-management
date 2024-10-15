package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/model"
	"github.com/gorilla/mux"
)

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
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Warehouse updated successfully")
}
