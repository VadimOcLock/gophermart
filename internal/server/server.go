package server

import (
	"github.com/VadimOcLock/gophermart/internal/accrualclient"
	"github.com/VadimOcLock/gophermart/internal/handler/balancehandler"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/service/balanceservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/balanceusecase"
	"github.com/VadimOcLock/gophermart/pkg/httpmix"

	"github.com/VadimOcLock/gophermart/internal/handler/orderhandler"
	"github.com/VadimOcLock/gophermart/internal/service/orderservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/orderusecase"

	"github.com/VadimOcLock/gophermart/internal/config"
	"github.com/VadimOcLock/gophermart/internal/handler/authhandler"
	"github.com/VadimOcLock/gophermart/internal/pgstore"
	"github.com/VadimOcLock/gophermart/internal/service/authservice"
	"github.com/VadimOcLock/gophermart/internal/usecase/authusecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/VadimOcLock/gophermart/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// New godoc
// @title			GopherMart API
// @version		v0.0.1
// @description	API накопительной системы лояльности.
// @contact.name	Kozenkov Vadim
// @contact.url	https://github.com/VadimOcLock
// @contact.email	vadimocloc@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
func New(pgClient *pgxpool.Pool, cfg config.WebServer) *chi.Mux {
	// Accrual client.
	accrualClient := accrualclient.NewAccrualClient(cfg.AccrualConfig.SrvAddr)
	// Store.
	store := pgstore.NewPgStore(pgClient)
	// Services.
	authService := authservice.NewAuthService(store)
	orderService := orderservice.NewOrderService(store, accrualClient)
	balanceService := balanceservice.NewBalanceService(store)
	// UseCases.
	authUseCase := authusecase.NewAuthUseCase(authService, authusecase.JWTConfig{
		SecretKey:     cfg.JWTConfig.SecretKey,
		TokenDuration: cfg.JWTConfig.TokenExpiration,
	})
	orderUseCase := orderusecase.NewOrderUseCase(orderService)
	balanceUseCase := balanceusecase.NewBalanceUseCase(balanceService)
	// Handler.
	authHandler := authhandler.New(authUseCase)
	orderHandler := orderhandler.New(orderUseCase)
	balanceHandler := balancehandler.New(balanceUseCase)
	// Server.
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Swagger.
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://"+cfg.WebServerConfig.SrvAddr+"/swagger/doc.json"),
	))

	// Auth.
	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)

	// Protected.
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(cfg.JWTConfig.SecretKey))
		r.Route("/user", func(r chi.Router) {
			r.Post("/orders", orderHandler.UploadOrder)
			r.Get("/orders", orderHandler.GetOrders)

			r.Get("/balance", balanceHandler.GetBalance)
			r.Post("/balance/withdraw", balanceHandler.WithdrawBalance)
			r.Get("/withdrawals", balanceHandler.GetWithdrawals)
		})
	})

	mux := httpmix.NewDefaultMux()
	mux.Mount("/", r)

	return mux
}
