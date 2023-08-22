package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "simple-url-shortener/internal/lib/api/response"
	"simple-url-shortener/internal/lib/logger/sl"
	"simple-url-shortener/internal/storage"
)

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	resp.Response
}

// URLDeleter is an interface for deleting url by alias
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
		}

		err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("alias not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to remove alias", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("removed alias", slog.String("alias", alias))

		responseOK(w, r, alias)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
