package main

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/service/auth"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/service/user"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/types"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  *user.Store
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Post("/login", app.loginHandler)
		r.Post("/logout", app.logoutHandler)
		r.Post("/timestamps", app.createTimestampHandler)
		r.Get("/timestamps", app.getAllTimestampsHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := app.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
	// Clear the JWT token or invalidate the session
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (app *application) createTimestampHandler(w http.ResponseWriter, r *http.Request) {
	var timestamp types.TimeStamp
	if err := utils.ParseJSON(r, &timestamp); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	timestamp.Year = int16(time.Now().Year())
	timestamp.Month = uint8(time.Now().Month())
	timestamp.Day = uint8(time.Now().Day())
	timestamp.Hour = uint8(time.Now().Hour())
	timestamp.Minute = uint8(time.Now().Minute())
	timestamp.Second = uint8(time.Now().Second())

	if err := app.store.CreateTimestamp(timestamp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, timestamp)
}

func (app *application) getAllTimestampsHandler(w http.ResponseWriter, r *http.Request) {
	timestamps, err := app.store.GetAllTimestamps()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, timestamps)
}
