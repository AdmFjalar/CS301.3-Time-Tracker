package main

import "net/http"

// healthCheckHandler is a simple handler function that responds with "ok".
// It is used to check the health status of the application.
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
