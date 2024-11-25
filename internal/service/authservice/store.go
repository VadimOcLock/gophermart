package authservice

import (
	"context"
	"time"
)

type UserStore interface {
	CreateUser(ctx context.Context, params CreateUserParams) (uint64, error)
	UserExistsByLogin(ctx context.Context, login string) (bool, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (uint64, error)
}

type CreateUserParams struct {
	Login        string
	PasswordHash string
}

type CreateSessionParams struct {
	UserID    uint64
	Token     string
	ExpiresAt time.Time
}
