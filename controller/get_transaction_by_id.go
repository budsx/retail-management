package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
