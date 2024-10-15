package controller

import "net/http"

func (c *Controller) Readiness(w http.ResponseWriter, r *http.Request) {
	err := c.service.Readiness(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Not Ready")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "ready")
}
