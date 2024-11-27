package main

import (
	"net/http"
)

// internalServerError godoc
//
//	@Summary		Internal Server Error
//	@Description	Logs an internal server error and writes a JSON error response with status 500
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/internal [post]
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

// forbiddenResponse godoc
//
//	@Summary		Forbidden
//	@Description	Logs a forbidden error and writes a JSON error response with status 403
//	@Tags			errors
//	@Produce		json
//	@Router			/errors/forbidden [post]
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

// badRequestResponse godoc
//
//	@Summary		Bad Request
//	@Description	Logs a bad request error and writes a JSON error response with status 400
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/bad-request [post]
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

// conflictResponse godoc
//
//	@Summary		Conflict
//	@Description	Logs a conflict error and writes a JSON error response with status 409
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/conflict [post]
func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusConflict, err.Error())
}

// notFoundResponse godoc
//
//	@Summary		Not Found
//	@Description	Logs a not found error and writes a JSON error response with status 404
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/not-found [post]
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

// unauthorizedErrorResponse godoc
//
//	@Summary		Unauthorized
//	@Description	Logs an unauthorized error and writes a JSON error response with status 401
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/unauthorized [post]
func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// unauthorizedBasicErrorResponse godoc
//
//	@Summary		Unauthorized Basic
//	@Description	Logs an unauthorized basic error, sets the WWW-Authenticate header, and writes a JSON error response with status 401
//	@Tags			errors
//	@Produce		json
//	@Param			error	body		error	true	"Error"
//	@Router			/errors/unauthorized-basic [post]
func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// rateLimitExceededResponse godoc
//
//	@Summary		Rate Limit Exceeded
//	@Description	Logs a rate limit exceeded error, sets the Retry-After header, and writes a JSON error response with status 429
//	@Tags			errors
//	@Produce		json
//	@Param			retryAfter	body		string	true	"Retry After"
//	@Router			/errors/rate-limit-exceeded [post]
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
