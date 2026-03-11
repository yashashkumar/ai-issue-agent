package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// fallback
		}
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string, requestID string) {
	writeJSON(w, status, ErrorResponse{Error: message, RequestID: requestID})
}
