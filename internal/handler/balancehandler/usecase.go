package balancehandler

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceUseCase interface {
	FindBalance(ctx context.Context, userID uint64) (entity.Balance, error)
}
