package helper

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func BindJson(r io.Reader, obj any) error {
	return json.NewDecoder(r).Decode(obj)
}
