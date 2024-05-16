package routes

import (
	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func Init(urls storage.Urls) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.WithRequestLogging)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.AddHandler(urls))
		r.Get("/{id}", handlers.FetchHandler(urls))
	})

	return r
}
