package orderservice

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/entity"
)

type OrderStore interface {
	FindOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
	SaveOrder(ctx context.Context, userID uint64, orderNumber string, status entity.OrderStatus) (uint64, error)
	UpdateOrderStatus(ctx context.Context, orderNumber string, status entity.OrderStatus) (bool, error)
	UpdateOrder(ctx context.Context, orderNumber string, status entity.OrderStatus, accrual float64) (bool, error)
	FindAllOrders(ctx context.Context, userID uint64) ([]entity.Order, error)
}
