package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// writeJSON godoc
//
//	@Summary		Writes JSON response
//	@Description	Writes the given data as a JSON response with the specified status code
//	@Tags			json
//	@Accept			json
//	@Produce		json
//	@Param			status	query		int		true	"Status code"
//	@Param			data	body		any		true	"Data to write"
//	@Success		200		{object}	any		"JSON response"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/json [post]
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// readJSON godoc
//
//	@Summary		Reads JSON request
//	@Description	Reads JSON data from the request body and decodes it into the provided data structure
//	@Tags			json
//	@Accept			json
//	@Produce		json
//	@Param			data	body		any		true	"Data to read"
//	@Success		200		{object}	any		"Decoded data"
//	@Failure		400		{object}	error	"Bad request"
//	@Router			/json [get]
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

// writeJSONError godoc
//
//	@Summary		Writes JSON error response
//	@Description	Writes a JSON error response with the specified status code and error message
//	@Tags			json
//	@Accept			json
//	@Produce		json
//	@Param			status	query		int		true	"Status code"
//	@Param			message	query		string	true	"Error message"
//	@Success		200		{object}	any		"JSON error response"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/json/error [post]
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

// jsonResponse godoc
//
//	@Summary		Writes JSON response
//	@Description	Writes a JSON response with the specified status code and data
//	@Tags			json
//	@Accept			json
//	@Produce		json
//	@Param			status	query		int		true	"Status code"
//	@Param			data	body		any		true	"Data to write"
//	@Success		200		{object}	any		"JSON response"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/json/response [post]
func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
