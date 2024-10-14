package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Controller) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := c.service.GetProductByID(r.Context(), productID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, product)
}
