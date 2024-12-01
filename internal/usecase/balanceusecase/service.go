package balanceusecase

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceService interface {
	FindBalance(ctx context.Context, userID uint64) (entity.Balance, error)
	FindWithdrawals(ctx context.Context, userID uint64) ([]entity.Withdraw, error)
}
