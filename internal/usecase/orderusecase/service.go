package orderusecase

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

type OrderService interface {
	UploadOrder(ctx context.Context, userID uint64, orderNumber string) (uint64, error)
	FindAllOrders(ctx context.Context, userID uint64) ([]entity.Order, error)
}
