package server

import (
	"net/http"
	"time"

	"github.com/VadimOcLock/gophermart/internal/config"
	"github.com/VadimOcLock/gophermart/internal/handler/authhandler"
	"github.com/VadimOcLock/gophermart/internal/pgstore"
	"github.com/VadimOcLock/gophermart/internal/service/authservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/authusecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pgClient *pgxpool.Pool, cfg config.WebServer) *http.Server {
	// Store.
	store := pgstore.New(pgClient)
	// Services.
	authService := authservice.NewAuthService(store)
	// UseCases.
	authUseCase := authusecase.NewAuthUseCase(authService, authusecase.JWTConfig{
		SecretKey:     cfg.JWTConfig.SecretKey,
		TokenDuration: cfg.JWTConfig.TokenExpiration,
	})
	// Handler.
	authHandler := authhandler.New(authUseCase)
	// Server.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	return &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
}
