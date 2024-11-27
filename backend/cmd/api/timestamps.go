package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
	"github.com/go-chi/chi/v5"
)

type timestampKey string

const timestampCtx timestampKey = "timestamp"

// CreateTimestampPayload represents the payload for creating a new timestamp.
type CreateTimestampPayload struct {
	StampType string `json:"stamp_type" validate:"required"`
	StampTime string `json:"stamp_time" validate:"required"`
}

// createTimestampHandler godoc
//
//	@Summary		Creates a timestamp
//	@Description	Creates a timestamp for a user
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTimestampPayload	true	"Timestamp information"
//	@Success		201		{object}	store.Timestamp			"Timestamp created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/timestamps [post]
func (app *application) createTimestampHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateTimestampPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//user := getUserFromContext(r)

	parsedTime, err := time.ParseInLocation(time.RFC3339, payload.StampTime, time.Local)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timestamp := &store.Timestamp{
		StampType: payload.StampType,
		UserID:    getUserFromContext(r).ID,
		StampTime: parsedTime,
	}

	ctx := r.Context()

	if err := app.store.Timestamps.Create(ctx, timestamp); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, timestamp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getTimestampHandler godoc
//
//	@Summary		Fetches a timestamp
//	@Description	Fetches a timestamp by ID
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Timestamp ID"
//	@Success		200	{object}	store.Timestamp
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/timestamps/{id} [get]
func (app *application) getTimestampHandler(w http.ResponseWriter, r *http.Request) {
	timestamp := getTimestampFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, timestamp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getLatestTimestampHandler godoc
//
//	@Summary		Fetches the latest timestamp
//	@Description	Fetches the most recent timestamp for a user
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.Timestamp
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/timestamps/latest [get]
func (app *application) getLatestTimestampHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	timestamp, err := app.store.Timestamps.GetLatestTimestamp(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, timestamp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getFinishedShiftsHandler godoc
//
//	@Summary		Fetches finished shifts
//	@Description	Fetches finished shifts for a user
//	@Tags			shifts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Shift
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/shifts [get]
func (app *application) getFinishedShiftsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	shifts, err := app.store.Timestamps.GetFinishedShifts(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, shifts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getFinishedShiftsByUserHandler godoc
//
//	@Summary		Fetches finished shifts by user ID
//	@Description	Fetches finished shifts for a specific user by their ID
//	@Tags			shifts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200	{object}	[]store.Shift
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/shifts/{userID} [get]
func (app *application) getFinishedShiftsByUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	idTemp, err := strconv.Atoi(idParam)
	id := int64(idTemp)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	shifts, err := app.store.Timestamps.GetFinishedShifts(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, shifts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// deleteTimestampHandler godoc
//
//	@Summary		Deletes a timestamp
//	@Description	Deletes a timestamp by ID
//	@Tags			timestamps
//	@Produce		json
//	@Param			id	path		int	true	"Timestamp ID"
//	@Success		204	{string}	string	"Timestamp deleted"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/timestamps/{id} [delete]
func (app *application) deleteTimestampHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "timestampID")
	idTemp, err := strconv.Atoi(idParam)
	id := int64(idTemp)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Timestamps.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTimestampPayload represents the payload for updating an existing timestamp.
type UpdateTimestampPayload struct {
	StampType string `json:"stamp_type" validate:"required"`
	StampTime string `json:"stamp_time" validate:"required"`
}

// updateTimestampHandler godoc
//
//	@Summary		Updates a timestamp
//	@Description	Updates a timestamp by ID
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Timestamp ID"
//	@Param			payload	body		UpdateTimestampPayload	true	"Updated timestamp information"
//	@Success		200		{object}	store.Timestamp
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/timestamps/{id} [patch]
func (app *application) updateTimestampHandler(w http.ResponseWriter, r *http.Request) {
	timestamp := getTimestampFromCtx(r)

	var payload UpdateTimestampPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, payload.StampTime)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timestamp.StampType = payload.StampType
	timestamp.StampTime = parsedTime

	ctx := r.Context()

	if err := app.updateTimestamp(ctx, timestamp); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, timestamp); err != nil {
		app.internalServerError(w, r, err)
	}
}

// timestampsContextMiddleware godoc
//
//	@Summary		Timestamps Context Middleware
//	@Description	Middleware that retrieves a timestamp by ID and adds it to the request context
//	@Tags			middleware
//	@Produce		json
//	@Router			/middleware/timestamps-context [get]
func (app *application) timestampsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "timestampID")
		idTemp, err := strconv.Atoi(idParam)
		id := int64(idTemp)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		timestamp, err := app.store.Timestamps.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, timestampCtx, timestamp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getTimestampFromCtx godoc
//
//	@Summary		Get Timestamp from Context
//	@Description	Retrieves the timestamp from the request context
//	@Tags			middleware
//	@Produce		json
//	@Router			/middleware/get-timestamp-from-ctx [get]
func getTimestampFromCtx(r *http.Request) *store.Timestamp {
	timestamp, _ := r.Context().Value(timestampCtx).(*store.Timestamp)
	return timestamp
}

// updateTimestamp godoc
//
//	@Summary		Update Timestamp
//	@Description	Updates a timestamp in the store and deletes the cache entry
//	@Tags			timestamps
//	@Produce		json
//	@Router			/timestamps/update [patch]
func (app *application) updateTimestamp(ctx context.Context, timestamp *store.Timestamp) error {
	if err := app.store.Timestamps.Update(ctx, timestamp); err != nil {
		return err
	}

	app.cacheStorage.Users.Delete(ctx, timestamp.ID)
	return nil
}
