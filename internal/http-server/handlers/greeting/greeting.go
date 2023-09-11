package greeting

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"simple-url-shortener/internal/lib/logger/sl"
)

func New(log *slog.Logger, staticDir string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.greeting.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		indexFilePath := fmt.Sprintf("%s/index.html", staticDir)
		indexHTML, err := readHtmlFromFile(indexFilePath)
		if err != nil {
			log.Error("Failed to read index.html", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "Internal Server Error")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(indexHTML))
	}
}

func readHtmlFromFile(filePath string) (string, error) {
	bs, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	return string(bs), nil
}
