package pgstore

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

const findOrderByOrderNumber = `
select id, user_id, order_number, status, accrual, uploaded_at
from orders
where order_number = $1
limit 1;
`

func (q *Queries) FindOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, error) {
	row := q.db.QueryRow(ctx, findOrderByOrderNumber, orderNumber)
	var order entity.Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

const saveOrder = `
insert into orders (user_id, order_number, status)
values ($1, $2, $3)
returning id;
`

func (q *Queries) SaveOrder(ctx context.Context, userID uint64, orderNumber string, status entity.OrderStatus) (uint64, error) {
	row := q.db.QueryRow(ctx, saveOrder, userID, orderNumber, status)
	var orderId uint64
	err := row.Scan(&orderId)

	return orderId, err
}

const updateOrderStatus = `
update orders
set status = $1
where order_number = $2
returning true as updated;
`

func (q *Queries) UpdateOrderStatus(ctx context.Context, orderNumber string, status entity.OrderStatus) (bool, error) {
	row := q.db.QueryRow(ctx, updateOrderStatus, status, orderNumber)
	var updated bool
	err := row.Scan(&updated)

	return updated, err
}

const updateOrder = `
update orders
set status  = $1,
    accrual = $2
where order_number = $3
returning true as updated;
`

func (q *Queries) UpdateOrder(ctx context.Context, orderNumber string, status entity.OrderStatus, accrual float64) (bool, error) {
	row := q.db.QueryRow(ctx, updateOrder, status, accrual, orderNumber)
	var updated bool
	err := row.Scan(&updated)

	return updated, err
}
