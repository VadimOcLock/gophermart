package balanceusecase

import (
	"context"
	"encoding/json"
	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/service/balanceservice"
)

type BalanceUseCase struct {
	BalanceService BalanceService
}

var _ BalanceService = (*balanceservice.BalanceService)(nil)

func NewBalanceUseCase(b BalanceService) *BalanceUseCase {
	return &BalanceUseCase{
		BalanceService: b,
	}
}

func (uc BalanceUseCase) FindBalance(ctx context.Context, userID uint64) (entity.Balance, error) {
	return uc.BalanceService.FindBalance(ctx, userID)
}

func (uc BalanceUseCase) FindWithdrawals(ctx context.Context, userID uint64) ([]byte, error) {
	ws, err := uc.BalanceService.FindWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(ws) == 0 {
		return nil, errorz.ErrUserHasNoWithdrawals
	}

	return json.Marshal(ws)
}
