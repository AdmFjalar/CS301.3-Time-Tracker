package utils

import (
	"encoding/json"
	"net/http"
)

// ParseJSON parses JSON from an HTTP request
func ParseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
