package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/docs"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/auth"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/env"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/mailer"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/ratelimiter"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store/cache"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:3000")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.config.rateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// Operations
		r.Get("/health", app.healthCheckHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/", app.getUserHandler)
			r.Patch("/change-password", app.changePasswordHandler)
		})

		r.Route("/timestamps", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createTimestampHandler)
			r.Get("/", app.getTimestampHandler)
			r.Get("/latest", app.getLatestTimestampHandler)

			r.Route("/{timestampID}", func(r chi.Router) {
				r.Use(app.timestampsContextMiddleware)
				r.Get("/", app.checkTimestampOwnership("manager", app.getTimestampHandler))

				r.Patch("/", app.checkRolePrecedenceMiddleware("manager", app.updateTimestampHandler))
				r.Delete("/", app.checkRolePrecedenceMiddleware("manager", app.deleteTimestampHandler))
			})
		})

		r.Route("/shifts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/", app.getFinishedShiftsHandler)
			r.Get("/{userID}", app.checkRolePrecedenceMiddleware("manager", app.getFinishedShiftsByUserHandler))
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.checkRolePrecedenceMiddleware("manager", app.getUsersHandler))
			})
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Patch("/", app.checkRolePrecedenceMiddleware("manager", app.updateUserHandler))
				r.Get("/", app.checkRolePrecedenceMiddleware("manager", app.getUserHandler))
				r.Delete("/", app.checkRolePrecedenceMiddleware("manager", app.deleteUserHandler))
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
