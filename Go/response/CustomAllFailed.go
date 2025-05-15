package response

import (
	"encoding/json"
	"net/http"
	"time"
)

func CustomAllFailed(w http.ResponseWriter) {
	errorResponse := map[string]interface{}{
		"status":    503,
		"message":   "All backends can not access. Please try later.",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(errorResponse)
}
