package controller

import (
	"net/http"
)

func (c *Controller) GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := c.service.GetStockTransactions(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendSuccessResponse(w, http.StatusOK, transactions)
}
