package server

import (
	"github.com/VadimOcLock/gophermart/cmd/external"
	"github.com/VadimOcLock/gophermart/internal/handler/orderhandler"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/service/orderservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/orderusecase"
	"net/http"
	"time"

	"github.com/VadimOcLock/gophermart/internal/config"
	"github.com/VadimOcLock/gophermart/internal/handler/authhandler"
	"github.com/VadimOcLock/gophermart/internal/pgstore"
	"github.com/VadimOcLock/gophermart/internal/service/authservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/authusecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pgClient *pgxpool.Pool, cfg config.WebServer) *http.Server {
	// Accrual client.
	accrualClient := external.NewAccrualClient(cfg.AccrualConfig.SrvAddr)
	// Store.
	store := pgstore.New(pgClient)
	// Services.
	authService := authservice.NewAuthService(store)
	orderService := orderservice.NewOrderService(store, accrualClient)
	// UseCases.
	authUseCase := authusecase.NewAuthUseCase(authService, authusecase.JWTConfig{
		SecretKey:     cfg.JWTConfig.SecretKey,
		TokenDuration: cfg.JWTConfig.TokenExpiration,
	})
	orderUseCase := orderusecase.NewOrderUseCase(orderService)
	// Handler.
	authHandler := authhandler.New(authUseCase)
	orderHandler := orderhandler.New(orderUseCase)
	// Server.
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Auth.
	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)

	// Protected.
	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(cfg.JWTConfig.SecretKey))
		r.Get("/protected", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Authorized endpoint"))
		})
		r.Post("/orders", orderHandler.UploadOrder)
	})

	return &http.Server{
		Addr:              cfg.WebServerConfig.SrvAddr,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
}
