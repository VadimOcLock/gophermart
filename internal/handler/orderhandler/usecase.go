package orderhandler

import (
	"context"
)

type OrderUseCase interface {
	UploadOrder(ctx context.Context, userID uint64, orderNumber string) error
	FindAllOrders(ctx context.Context, userID uint64) ([]byte, error)
}
