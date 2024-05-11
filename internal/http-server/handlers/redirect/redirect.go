package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"urlshort/internal/lib/api/response"
	"urlshort/internal/lib/logger"
	"urlshort/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.redirect.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias") //getting url params
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("not found"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", logger.Err(err))
			render.JSON(w, r, response.Error("not found"))
			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)

	}
}
