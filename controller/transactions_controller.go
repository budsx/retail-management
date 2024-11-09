package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
	"github.com/gorilla/mux"
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

func (c *Controller) GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := c.service.GetStockTransactions(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, transactions)
}

func (c *Controller) GetStockTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	transactionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	transaction, err := c.service.GetStockTransactionByID(r.Context(), transactionID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, transaction)
}

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
