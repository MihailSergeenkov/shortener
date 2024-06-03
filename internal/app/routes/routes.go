package routes

import (
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(l *zap.Logger, s data.Storager) chi.Router {
	r := chi.NewRouter()
	r.Use(withRequestLogging(l), gzipMiddleware(l))

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.PingHandler(l, s))
		r.Post("/", handlers.AddHandler(l, s))
		r.Get("/{id}", handlers.FetchHandler(l, s))

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Route("/api", func(r chi.Router) {
				r.Post("/shorten", handlers.APIAddHandler(l, s))
			})
		})
	})

	return r
}
