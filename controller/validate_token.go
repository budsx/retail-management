package controller

import "net/http"

func (c *Controller) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value("user_id").(int64)
	username := r.Context().Value("username").(string)

	sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"username": username,
		"valid":    true,
	})
}
