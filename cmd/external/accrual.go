package external

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const defaultClientTimeout = 5 * time.Second

type AccrualClient struct {
	BaseURL string
	Client  *resty.Client
}

func NewAccrualClient(baseURL string) *AccrualClient {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(defaultClientTimeout).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		})

	return &AccrualClient{
		BaseURL: baseURL,
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
	resp, err := c.Client.R().
		SetContext(ctx).
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
