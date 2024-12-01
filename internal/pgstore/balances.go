package pgstore

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
)

const findBalanceByUserID = `
select (select coalesce(sum(accrual), 0) from orders where user_id = $1) -
       (select coalesce(sum(sum), 0) from withdrawals where user_id = $1) as current_balance,
       (select coalesce(sum(sum), 0) from withdrawals where user_id = $1) as withdrawn_balance
`

func (q *Queries) FindBalanceByUserID(ctx context.Context, userID uint64) (entity.Balance, error) {
	row := q.db.QueryRow(ctx, findBalanceByUserID, userID)
	var balance entity.Balance
	err := row.Scan(
		&balance.CurrentBalance,
		&balance.WithdrawnBalance)

	return balance, err
}
