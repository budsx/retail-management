package controller

import "net/http"

func (c *Controller) GetTotalStocks(w http.ResponseWriter, r *http.Request) {

	totalStock, err := c.service.GetTotalStocks(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, totalStock)
}
