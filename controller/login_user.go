package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/model"
	"github.com/budsx/retail-management/utils"
)

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.service.ValidateUser(r.Context(), model.Credentials{
		Username: creds.Username,
		Password: creds.Password,
	})
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := utils.GenerateJWT(int64(user.UserID), user.Username)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	sendSuccessResponse(w, http.StatusOK, map[string]string{"token": token})
}
