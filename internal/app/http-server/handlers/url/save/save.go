package save

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	resp "simple-url-shortener/internal/app/lib/api/response"
	"simple-url-shortener/internal/app/lib/logger/sl"
	"simple-url-shortener/internal/app/lib/random"
	"simple-url-shortener/internal/app/repository"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	Add(ctx context.Context, Url *repository.Url) (int64, error)
	CheckAliasExists(ctx context.Context, alias string) (bool, error)
	CheckURLExists(ctx context.Context, urlToCheck string) (bool, error)
	GetAliasByURL(ctx context.Context, urlToFind string) (string, error)
}

func New(ctx context.Context, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode requested body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidateResponse(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {

			exists, err := urlSaver.CheckURLExists(ctx, req.URL)
			if err != nil {
				log.Error("failed to check that URL exists in DB", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to check that URL exists in DB"))
				return
			}

			if exists {
				alias, err = urlSaver.GetAliasByURL(ctx, req.URL)
				if err != nil {
					log.Error("failed to get alias connected to URL", sl.Err(err))
					render.JSON(w, r, resp.Error("failed to get alias connected to URL"))
					return
				}

				responseOK(w, r, alias)
				return
			}

			const maxAttempts = 64 // Maximum number of generation attempts
			exists = true

			for attempt := 1; attempt <= maxAttempts; attempt++ {
				alias = random.NewRandomString()
				exists, err = urlSaver.CheckAliasExists(ctx, alias)
				if err != nil {
					log.Error("failed to generate alias", sl.Err(err))
					render.JSON(w, r, resp.Error("failed to generate url"))
					return
				}
				if !exists {
					break
				}
			}

			if exists {
				log.Error("The number of attempts to create an alias has been exceeded", sl.Err(err))
				render.JSON(w, r, resp.Error("The number of attempts to create an alias has been exceeded. Try again after a while"))
				return
			}
		}

		insertUrl := repository.Url{
			Alias: alias,
			URL:   req.URL,
		}

		id, err := urlSaver.Add(ctx, &insertUrl)

		if errors.Is(err, repository.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
