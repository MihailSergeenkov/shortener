package routes

import (
	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
	api_handlers "github.com/MihailSergeenkov/shortener/internal/app/handlers/api"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init(urls storage.Urls) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithRequestLogging, gzipMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.AddHandler(urls))
		r.Get("/{id}", handlers.FetchHandler(urls))

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Route("/api", func(r chi.Router) {
				r.Post("/shorten", api_handlers.APIAddHandler(urls))
			})
		})
	})

	return r
}
