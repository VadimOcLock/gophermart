package pgstore

import (
	"context"
	"errors"
	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/jackc/pgx/v5"
)

const findAllWithdrawalsByUserID = `
select id, user_id, order_number, sum, processed_at
from withdrawals
where user_id = $1
order by processed_at desc;
`

func (q *Queries) FindAllWithdrawalsByUserID(ctx context.Context, userID uint64) ([]entity.Withdraw, error) {
	rows, err := q.db.Query(ctx, findAllWithdrawalsByUserID, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()
	var ws []entity.Withdraw
	for rows.Next() {
		var w entity.Withdraw
		err = rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderNumber,
			&w.Sum,
			&w.ProcessedAt)
		ws = append(ws, w)
	}

	return ws, nil
}

const withdrawal = `
insert into withdrawals(user_id, order_number, sum)
values ($1, $2, $3)
returning id;
`

func (s *PgStore) Withdrawal(ctx context.Context, userID uint64, orderNumber string, sum float64) (uint64, error) {
	row := s.db.QueryRow(ctx, withdrawal, userID, orderNumber, sum)
	var id uint64
	err := row.Scan(&id)

	return id, err
}
