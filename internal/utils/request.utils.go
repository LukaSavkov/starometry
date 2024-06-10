package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONFromRequest(r *http.Request, rw http.ResponseWriter, v interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		WriteErrorResp("Internal server error", 500, r.URL.Path, rw)
		return false
	}
	return true
}
