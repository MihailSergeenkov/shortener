// Модель роутинга сервиса.
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
)

// NewRouter функция инициализации роутинга.
func NewRouter(l *zap.Logger, s data.Storager) chi.Router {
	r := chi.NewRouter()
	r.Use(withRequestLogging(l))
	r.Mount("/debug", middleware.Profiler())

	r.Get("/ping", handlers.PingHandler(l, s))

	r.Route("/", func(r chi.Router) {
		r.Use(setAuthMiddleware(l), gzipMiddleware(l))
		r.Post("/", handlers.AddHandler(l, s))
		r.Get("/{id}", handlers.FetchHandler(l, s))

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType(common.JSONContentType))

			r.Route("/api", func(r chi.Router) {
				r.Route("/shorten", func(r chi.Router) {
					r.Post("/", handlers.APIAddHandler(l, s))
					r.Post("/batch", handlers.APIAddBatchHandler(l, s))
				})
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType(common.JSONContentType), checkAuthMiddleware(l))

		r.Route("/api/user/urls", func(r chi.Router) {
			r.Get("/", handlers.APIFetchUserURLsHandler(l, s))
			r.Delete("/", handlers.APIDeleteUserURLsHandler(l, s))
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType(common.JSONContentType), checkSubnetMiddleware(l, config.Params.TrustedSubnet))

		r.Get("/api/internal/stats", handlers.APIFetchStatsHandler(l, s))
	})

	return r
}
