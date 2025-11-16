package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

func newErrorResponse(message string, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		_ = err
	}
}


func newSuccessResponse(data interface{}, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		_ = err
		newErrorResponse("Internal server error", w, http.StatusInternalServerError)
	}
}
