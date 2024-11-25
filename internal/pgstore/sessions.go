package pgstore

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/service/authservice"
)

const createSession = `
INSERT INTO sessions (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING id;
`

func (q *Queries) CreateSession(ctx context.Context, params authservice.CreateSessionParams) (uint64, error) {
	row := q.db.QueryRow(ctx, createSession, params.UserID, params.Token, params.ExpiresAt)
	var id uint64
	err := row.Scan(&id)

	return id, err
}
