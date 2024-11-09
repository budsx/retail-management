package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
	"github.com/budsx/retail-management/utils"
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
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "User registered successfully")
}

func (c *Controller) ValidateToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextKeyUserID).(int64)
	if userID == 0 {
		sendErrorResponse(w, http.StatusUnauthorized, "Unathorized")
	}
	username := r.Context().Value(middleware.ContextKeyUsername).(string)
	if username == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "Unathorized")
	}

	sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"user_id":  userID,
		"username": username,
	})
}

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
