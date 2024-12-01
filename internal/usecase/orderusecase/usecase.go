package orderusecase

import (
	"context"
	"encoding/json"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/service/orderservice"
)

type OrderUseCase struct {
	OrderService OrderService
}

var _ OrderService = (*orderservice.OrderService)(nil)

func NewOrderUseCase(orderService OrderService) *OrderUseCase {
	return &OrderUseCase{OrderService: orderService}
}

func (uc OrderUseCase) UploadOrder(ctx context.Context, userID uint64, orderNumber string) error {
	if !IsValidOrderNumber(orderNumber) {
		return errorz.ErrInvalidOrderNumberFormat
	}
	_, err := uc.OrderService.UploadOrder(ctx, userID, orderNumber)
	if err != nil {
		return err
	}

	return nil
}

func IsValidOrderNumber(orderNumber string) bool {
	var sum int
	var double bool

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit := orderNumber[i]
		if digit < '0' || digit > '9' {
			return false
		}

		n := int(digit - '0')
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		double = !double
	}

	return sum%10 == 0
}

func (uc OrderUseCase) FindAllOrders(ctx context.Context, userID uint64) ([]byte, error) {
	orders, err := uc.OrderService.FindAllOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errorz.ErrUserHasNoOrders
	}

	return json.Marshal(orders)
}
