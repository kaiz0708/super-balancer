package response

import (
	"encoding/json"
	"net/http"
	"time"
)

func CustomAllFailed(w http.ResponseWriter) {
	errorResponse := map[string]interface{}{
		"status":    503,
		"message":   "Tất cả backend hiện tại đều không thể truy cập. Xin thử lại sau.",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(errorResponse)
}
