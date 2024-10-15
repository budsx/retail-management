package controller

import "net/http"

func (c *Controller) GetWarehousesByUserID(w http.ResponseWriter, r *http.Request) {
	warehouses, err := c.service.GetWarehouseByUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, warehouses)
}
