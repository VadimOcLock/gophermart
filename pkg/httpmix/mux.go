package httpmix

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMux(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middlewares...)

	return r
}

func NewDefaultMux() *chi.Mux {
	return NewMux(middleware.RequestID, middleware.RealIP, middleware.Logger, Cors)
}
