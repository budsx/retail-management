package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/model"
	"github.com/gorilla/mux"
)

func (c *Controller) EditProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Missing product ID")
		return
	}

	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var updatedProduct model.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedProduct.ProductID = productID
	err = c.service.EditProduct(r.Context(), updatedProduct)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	sendSuccessResponse(w, http.StatusOK, "Product updated successfully")
}
