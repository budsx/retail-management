package controller

import (
	"net/http"
	"strconv"

	"github.com/budsx/retail-management/model"
)

func (c *Controller) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	pagination := model.Pagination{
		Page:  int32(page),
		Limit: int32(limit),
	}

	products, err := c.service.GetProducts(r.Context(), pagination)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	sendSuccessResponse(w, http.StatusOK, products)
}
