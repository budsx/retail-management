package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Controller) GetTotalStockByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationIDStr := vars["location_id"]
	locationID, err := strconv.ParseInt(locationIDStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	totalStock, err := c.service.GetTotalStockByLocation(r.Context(), locationID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, totalStock)
}
