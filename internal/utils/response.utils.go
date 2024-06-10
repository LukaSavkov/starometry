package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/c12s/metrics/internal/models"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Handle JSON encoding error
		WriteErrorResp(err.Error(), http.StatusInternalServerError, "Internal Server Error", w)
	}
}

func WriteErrorResp(err string, status int, path string, w http.ResponseWriter) {

	baseErrorResp := models.BaseErrorHttpResponse{
		Error:  err,
		Path:   path,
		Status: status,
		Time:   time.Now().String(),
	}
	writeJSONResponse(w, status, baseErrorResp)
}

func WriteResp(resp any, statusCode int, w http.ResponseWriter) {
	if resp == nil {
		return
	}
	domainResponse := models.BaseHttpResponse{
		Status: statusCode,
		Data:   resp,
	}
	writeJSONResponse(w, statusCode, domainResponse)
}
