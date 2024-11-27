package orderusecase

import "context"

type OrderService interface {
	UploadOrder(ctx context.Context, userID uint64, orderNumber string) (uint64, error)
}
