package orderservice

import (
	"context"
	"github.com/VadimOcLock/gophermart/cmd/external"
	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/rs/zerolog/log"
)

type OrderService struct {
	OrderStore    OrderStore
	AccrualClient *external.AccrualClient
}

func NewOrderService(orderStore OrderStore, accrualClient *external.AccrualClient) *OrderService {
	return &OrderService{
		OrderStore:    orderStore,
		AccrualClient: accrualClient,
	}
}

func (s OrderService) UploadOrder(ctx context.Context, userID uint64, orderNumber string) (uint64, error) {
	order, err := s.OrderStore.FindOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		return 0, err
	}
	if order != nil {
		if order.UserID == userID {
			return 0, errorz.ErrOrderAlreadyUploadedByUser
		}

		return 0, errorz.ErrOrderAlreadyUploadedByAnotherUser
	}
	orderID, err := s.OrderStore.SaveOrder(ctx, userID, orderNumber, entity.OrderStatusNew)
	if err != nil {
		return 0, err
	}
	go s.processOrder(ctx, orderNumber)

	return orderID, nil
}

func (s OrderService) processOrder(ctx context.Context, orderNumber string) {
	// PROCESSING
	if err := s.OrderStore.UpdateOrderStatus(ctx, orderNumber, entity.OrderStatusProcessing); err != nil {
		log.Error().Err(err).Msg("failed to update order status")

		return
	}
	resp, err := s.AccrualClient.GetOrderAccrual(ctx, orderNumber)
	if err != nil {
		log.Error().Msg(err.Error())

		return
	}
	if resp == nil {
		// No content.

		return
	}
	switch resp.Status {
	case string(entity.OrderStatusProcessed):
		// todo
	}

}
