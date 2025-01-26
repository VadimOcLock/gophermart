package balanceservice

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
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

func (s BalanceService) OrderNumberExists(ctx context.Context, userID uint64, orderNumber string) (bool, error) {
	return s.BalanceStore.OrderNumberExists(ctx, userID, orderNumber)
}

func (s BalanceService) Withdrawal(
	ctx context.Context,
	userID uint64,
	orderNumber string,
	sum float64,
) (uint64, error) {
	balance, err := s.BalanceStore.FindBalanceByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if balance.CurrentBalance < sum {
		return 0, errorz.ErrNotEnoughFundsOnBalance
	}

	return s.BalanceStore.Withdrawal(ctx, userID, orderNumber, sum)
}
