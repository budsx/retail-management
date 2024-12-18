package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/model"
	"github.com/gorilla/mux"
)

func (c *Controller) AddLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.service.AddLocation(r.Context(), location)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Location added successfully")
}

func (c *Controller) EditLocationByUserID(w http.ResponseWriter, r *http.Request) {
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

	var location model.Location
	err = json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	location.LocationID = locationID
	err = c.service.EditLocationByUserID(r.Context(), location)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Location updated successfully")
}
