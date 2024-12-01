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

func (s BalanceService) FindWithdrawals(ctx context.Context, userID uint64) ([]entity.Withdraw, error) {
	return s.BalanceStore.FindAllWithdrawalsByUserID(ctx, userID)
}
