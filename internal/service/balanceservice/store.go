package balanceservice

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceStore interface {
	FindBalanceByUserID(ctx context.Context, userID uint64) (entity.Balance, error)
}
