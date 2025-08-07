// Package server provides the HTTP server setup and lifecycle management.
package server

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mexirica/chi-template/internal/configs"
	"github.com/mexirica/chi-template/internal/db/repository"
	"github.com/mexirica/chi-template/internal/handler"
	"github.com/mexirica/chi-template/internal/helpers"
	"github.com/mexirica/chi-template/internal/middleware"
	"github.com/mexirica/chi-template/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/mexirica/chi-template/docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	stdmiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	cfg         *configs.Config
	redis       *redis.Client
	db          *pgxpool.Pool
	srv         *http.Server
	userHandler *handler.MovieHandler
}

func New(cfg *configs.Config, redis *redis.Client, db *pgxpool.Pool) *App {
	userRepo := repository.NewMovieRepository(db)
	userService := service.NewMovieService(userRepo)
	userHandler := handler.NewMovieHandler(userService)

	app := &App{
		cfg:         cfg,
		redis:       redis,
		db:          db,
		userHandler: userHandler,
	}

	app.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.PORT),
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return app
}

func (app *App) Serve() {
	log.Info().Msgf("Starting server on port %s", app.cfg.PORT)
	if err := app.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Msgf("Erro inesperado no servidor: %v", err)
	}
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.srv.Shutdown(ctx)
}

func (app *App) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Logger)
	r.Use(std.HandlerProvider("", stdmiddleware.New(stdmiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, "API is up and running")
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", app.cfg.PORT))),
	)

	r.Route("/movies", func(r chi.Router) {
		r.Post("/", app.userHandler.Create)
		r.Delete("/{id}", app.userHandler.Delete)
		r.With(middleware.CacheMiddleware(time.Hour*24)).Get("/{id}", app.userHandler.GetById)
		r.With(middleware.CacheMiddleware(time.Hour*24)).Get("/list", app.userHandler.GetList)
	})

	return r
}
