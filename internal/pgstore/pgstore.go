package pgstore

import (
	"context"
	"fmt"
	"github.com/VadimOcLock/gophermart/internal/service/orderservice"

	"github.com/VadimOcLock/gophermart/internal/service/authservice"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewPgStore(db *pgxpool.Pool) *PgStore {
	return &PgStore{
		Queries: New(db),
		db:      db,
	}
}

type Store interface {
	authservice.UserStore
	orderservice.OrderStore
}

var _ Store = (*PgStore)(nil)

func (s *PgStore) ExecTx(ctx context.Context, fn func(queries *Queries) error) error {
	var tx, err = s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := s.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v; rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
