package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
)

func (c *Controller) CreateStockTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextKeyUserID).(int64)
	if userID == 0 {
		sendErrorResponse(w, http.StatusUnauthorized, "Unathorized")
	}

	var stockTransaction model.StockTransaction
	err := json.NewDecoder(r.Body).Decode(&stockTransaction)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	stockTransaction.CreatedBy = userID
	err = c.service.CreateStockTransaction(r.Context(), stockTransaction)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "Stock transaction created successfully")
}
