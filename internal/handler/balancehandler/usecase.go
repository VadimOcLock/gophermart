package balancehandler

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceUseCase interface {
	FindBalance(ctx context.Context, userID uint64) (entity.Balance, error)
	FindWithdrawals(ctx context.Context, userID uint64) ([]byte, error)
	Withdrawal(ctx context.Context, userID uint64, sum float64, orderNumber string) error
}
