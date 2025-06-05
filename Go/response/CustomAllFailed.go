package response

import (
	"encoding/json"
	"net/http"
	"time"
)

func CustomAllFailed(w http.ResponseWriter) {
	errorResponse := map[string]interface{}{
		"status":    503,
		"message":   "All backends can not access",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse)
}
