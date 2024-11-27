package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map to hold the health check data
	data := map[string]string{
		"status":  "ok",          // Status of the application
		"env":     app.config.env, // Current environment (e.g., development, production)
		"version": version,        // Application version
	}

	// Send the health check data as a JSON response
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		// Handle any errors that occur while sending the response
		app.internalServerError(w, r, err)
	}
}
