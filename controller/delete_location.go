package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Controller) DeleteLocationByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Missing location ID")
		return
	}

	locationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	err = c.service.DeleteLocationByUserID(r.Context(), locationID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Location deleted successfully")
}
