package utils

import (
	"encoding/json"
	"net/http"
)

// ParseJSON parses JSON from an HTTP request
// It decodes the JSON request body into the provided interface.
func ParseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// WriteError writes an error response
// It sets the HTTP status code and writes the error message as the response body.
func WriteError(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}

// WriteJSON writes a JSON response
// It sets the Content-Type header to "application/json" and writes the JSON-encoded response body.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
