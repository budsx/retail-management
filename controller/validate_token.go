package controller

import (
	"net/http"

	"github.com/budsx/retail-management/middleware"
)

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
		"status":   "valid",
	})
}
