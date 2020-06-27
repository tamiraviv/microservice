package rest

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Adapter) newRouter(timeout time.Duration) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(timeout))
	r.Route("/documents", func(r chi.Router) {
		r.Get("/{id}", s.getDocument)
		r.Post("/", s.addDocument)
	})
	return r
}
