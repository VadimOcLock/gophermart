package balanceservice

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type BalanceService struct {
	BalanceStore BalanceStore
}

func NewBalanceService(balanceStore BalanceStore) *BalanceService {
	return &BalanceService{
		BalanceStore: balanceStore,
	}
}

func (s BalanceService) FindBalance(ctx context.Context, userID uint64) (entity.Balance, error) {
	return s.BalanceStore.FindBalanceByUserID(ctx, userID)
}
