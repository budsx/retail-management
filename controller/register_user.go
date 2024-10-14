package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/model"
)

func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = c.service.RegisterUser(r.Context(), user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "User registered successfully")
}
