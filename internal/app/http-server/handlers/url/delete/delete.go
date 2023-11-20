package delete

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "simple-url-shortener/internal/app/lib/api/response"
	"simple-url-shortener/internal/app/lib/logger/sl"
	"simple-url-shortener/internal/app/repository"
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
	DeleteByAlias(ctx context.Context, alias string) error
}

func New(ctx context.Context, log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
		}

		err := urlDeleter.DeleteByAlias(ctx, alias)

		if errors.Is(err, repository.ErrURLNotFound) {
			log.Info("alias not found", "alias", alias)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to remove alias", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("removed alias", slog.String("alias", alias))

		responseOK(w, r)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
