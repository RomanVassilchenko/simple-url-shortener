package getAll

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
	resp "simple-url-shortener/internal/app/lib/api/response"
	"simple-url-shortener/internal/app/lib/logger/sl"
	"simple-url-shortener/internal/app/repository"
)

type Response struct {
	resp.Response
	Urls *[]repository.Url `json:"urls,omitempty"`
}

// URLGetter is an interface for getting all urls
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetAll(ctx context.Context) (*[]repository.Url, error)
}

func New(ctx context.Context, log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetAll.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		resURLs, err := urlGetter.GetAll(ctx)

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got urls")

		responseOK(w, r, resURLs)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, urls *[]repository.Url) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Urls:     urls,
	})
}
