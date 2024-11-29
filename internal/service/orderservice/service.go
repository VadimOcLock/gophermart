package orderservice

import (
	"context"
	"errors"
	"github.com/VadimOcLock/gophermart/cmd/external"
	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/jackc/pgx/v5"
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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
	go s.processOrder(orderNumber)

	return orderID, nil
}

func (s OrderService) processOrder(orderNumber string) {
	ctx := context.Background()
	resp, err := s.AccrualClient.GetOrderAccrual(ctx, orderNumber)
	switch {
	case err != nil:
		log.Error().Msg(err.Error())
	// 204.
	case resp == nil:
		if _, err = s.OrderStore.UpdateOrderStatus(ctx, orderNumber, entity.OrderStatusInvalid); err != nil {
			log.Error().Err(err).Msg("failed to update order status with invalid status")
		}
	case resp.Status == string(external.OrderStatusProcessing):
		if _, err = s.OrderStore.UpdateOrderStatus(ctx, orderNumber, entity.OrderStatusProcessing); err != nil {
			log.Error().Err(err).Msg("failed to update order status with processing status")
		}
	case resp.Status == string(external.OrderStatusInvalid):
		if _, err = s.OrderStore.UpdateOrderStatus(ctx, orderNumber, entity.OrderStatusInvalid); err != nil {
			log.Error().Err(err).Msg("failed to update order status with invalid status")
		}
	case resp.Status == string(external.OrderStatusProcessed):
		var accrual float64
		if resp.Accrual != nil {
			accrual = *resp.Accrual
		}
		if _, err = s.OrderStore.UpdateOrder(ctx, orderNumber, entity.OrderStatusProcessed, accrual); err != nil {
			log.Error().Err(err).Msg("failed to update order with processed status")
		}
	case resp.Status == string(external.OrderStatusRegistered):
		return
	}
}
