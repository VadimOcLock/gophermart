package balanceusecase

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceService interface {
	FindBalance(ctx context.Context, userID uint64) (entity.Balance, error)
	FindWithdrawals(ctx context.Context, userID uint64) ([]entity.Withdraw, error)
	OrderNumberExists(ctx context.Context, userID uint64, orderNumber string) (bool, error)
	Withdrawal(ctx context.Context, userID uint64, orderNumber string, sum float64) (uint64, error)
}
