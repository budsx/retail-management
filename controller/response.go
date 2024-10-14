package controller

import (
	"encoding/json"
	"net/http"

	"github.com/budsx/retail-management/utils"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func sendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := SuccessResponse{
		Message: utils.SuccessMessage,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Error: errorMessage,
	}
	json.NewEncoder(w).Encode(response)
}
