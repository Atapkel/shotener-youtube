package save

import (
	"errors"
	"github.com/atapkel/shortener/internal/lib/api/response"
	"github.com/atapkel/shortener/internal/lib/logger/sl"
	"github.com/atapkel/shortener/internal/lib/random"
	"github.com/atapkel/shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 7

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("Request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := saver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Info("failed to add db", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add db"))
			return
		}

	}

}
