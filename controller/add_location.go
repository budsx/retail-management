package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/model"
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
