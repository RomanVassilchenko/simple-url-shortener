package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"simple-url-shortener/internal/app/config"
	"simple-url-shortener/internal/app/db"
	"simple-url-shortener/internal/app/http-server/handlers/greeting"
	"simple-url-shortener/internal/app/http-server/handlers/redirect"
	deleteRequest "simple-url-shortener/internal/app/http-server/handlers/url/delete"
	"simple-url-shortener/internal/app/http-server/handlers/url/getAll"
	"simple-url-shortener/internal/app/http-server/handlers/url/save"
	mwLogger "simple-url-shortener/internal/app/http-server/middleware/logger"
	"simple-url-shortener/internal/app/lib/logger"
	"simple-url-shortener/internal/app/repository/postgresql"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env, cfg.LoggerPath)

	log.Info("starting simple-url-shortener", slog.String("env", cfg.Env))

	log.Debug("debug messages are enabled")

	database, err := db.NewDB(ctx, &cfg.DatabaseCredentials)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer database.GetPool(ctx).Close()

	urlsRepo := postgresql.NewUrls(database)

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	fsWeb := http.FileServer(http.Dir("./static/web"))
	router.Handle("/web/*", http.StripPrefix("/web/", fsWeb))

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("simple-url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Get("/getAll", getAll.New(ctx, log, urlsRepo))
		r.Delete("/{alias}", deleteRequest.New(ctx, log, urlsRepo))
	})

	router.Get("/", greeting.New(log, "./static"))
	router.Post("/", save.New(ctx, log, urlsRepo))
	router.Get("/{alias}", redirect.New(ctx, log, urlsRepo))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start Server", slog.String("error", err.Error()))
	}

	log.Error("server stopped")
}
