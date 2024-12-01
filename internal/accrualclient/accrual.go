package accrualclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const defaultClientTimeout = 5 * time.Second
const HTTPProtocol = "http://"

type AccrualClient struct {
	BaseURL string
	Client  *resty.Client
}

func NewAccrualClient(addr string) *AccrualClient {
	client := resty.New().
		SetBaseURL(HTTPProtocol + addr).
		SetTimeout(defaultClientTimeout).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if r.StatusCode() == http.StatusTooManyRequests ||
				r.StatusCode() == http.StatusNoContent || err != nil {
				return true
			}

			return false
		})

	return &AccrualClient{
		BaseURL: addr,
		Client:  client,
	}
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

func (c AccrualClient) GetOrderAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error) {
	retryConditionFn := func(r *resty.Response, err error) bool {
		var orderStatus string
		if r.Result().(*AccrualResponse) != nil {
			orderStatus = r.Result().(*AccrualResponse).Status
		}
		return r.StatusCode() == http.StatusNoContent || orderStatus == string(OrderStatusRegistered)
	}
	resp, err := c.Client.R().
		SetContext(ctx).
		AddRetryCondition(retryConditionFn).
		SetPathParam("orderNumber", orderNumber).
		SetResult(&AccrualResponse{}).
		Get("/api/orders/{orderNumber}")
	if err != nil {
		return nil, fmt.Errorf("GetOrderAccrual err: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Result().(*AccrualResponse), nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		return nil, errors.New("rate limit exceeded, try again later")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode())
	}
}
