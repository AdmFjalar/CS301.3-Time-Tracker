package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
	"github.com/go-chi/chi/v5"
)

type timestampKey string

const timestampCtx timestampKey = "timestamp"

type CreateTimestampPayload struct {
	StampType   string `json:"stamp_type"`
	UserID      int64  `json:"user_id"`
	TimeStampID int64  `json:"timestamp_id"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	Hour        int    `json:"hour"`
	Minute      int    `json:"minute"`
	Second      int    `json:"second"`
}

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

	user := getUserFromContext(r)

	timestamp := &store.Timestamp{
		StampType: payload.StampType,
		UserID:    user.ID,
		Year:      payload.Year,
		Month:     payload.Month,
		Day:       payload.Day,
		Hour:      payload.Hour,
		Minute:    payload.Minute,
		Second:    payload.Second,
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

func (app *application) getAllTimestampsHandler(w http.ResponseWriter, r *http.Request) {
	timestamp := getTimestampFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, timestamp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

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

type UpdateTimestampPayload struct {
	StampType string `json:"stamp_type"`
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	Second    int    `json:"second"`
}

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

	timestamp.StampType = payload.StampType
	timestamp.Year = payload.Year
	timestamp.Month = payload.Month
	timestamp.Day = payload.Day
	timestamp.Hour = payload.Hour
	timestamp.Minute = payload.Minute
	timestamp.Second = payload.Second

	ctx := r.Context()

	if err := app.updateTimestamp(ctx, &timestamp); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, timestamp); err != nil {
		app.internalServerError(w, r, err)
	}
}

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

func getTimestampFromCtx(r *http.Request) store.Timestamp {
	timestamp, _ := r.Context().Value(timestampCtx).(store.Timestamp)
	return timestamp
}

func (app *application) updateTimestamp(ctx context.Context, timestamp *store.Timestamp) error {
	if err := app.store.Timestamps.Update(ctx, timestamp); err != nil {
		return err
	}

	app.cacheStorage.Users.Delete(ctx, timestamp.ID)
	return nil
}
