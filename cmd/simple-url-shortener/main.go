package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"simple-url-shortener/internal/config"
	"simple-url-shortener/internal/http-server/handlers/greeting"
	"simple-url-shortener/internal/http-server/handlers/redirect"
	"simple-url-shortener/internal/http-server/handlers/url/delete"
	"simple-url-shortener/internal/http-server/handlers/url/save"
	mwLogger "simple-url-shortener/internal/http-server/middleware/logger"
	"simple-url-shortener/internal/lib/logger/handlers/slogpretty"
	"simple-url-shortener/internal/lib/logger/sl"
	"simple-url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log, err := setupLogger(cfg.Env, cfg.LoggerPath)

	if err != nil {
		os.Exit(1)
	}

	log.Info("starting simple-url-shortener", slog.String("env", cfg.Env))

	//log.Info("info messages are enabled")
	log.Debug("debug messages are enabled")
	//log.Error("error messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

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

		//r.Post("/", save.New(log, storage))
		//r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/", greeting.New(log, "./static"))
	router.Post("/", save.New(log, storage))
	router.Delete("/{alias}", delete.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start Server")
	}

	log.Error("server stopped")
}

func setupLogger(env string, logFilePath string) (*slog.Logger, error) {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log, nil
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
