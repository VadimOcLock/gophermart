package balanceusecase

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
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
