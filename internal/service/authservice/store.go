package authservice

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/entity"
	"time"
)

type UserStore interface {
	CreateUser(ctx context.Context, params CreateUserParams) (uint64, error)
	UserExistsByLogin(ctx context.Context, login string) (bool, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (uint64, error)
	FindUserByLogin(ctx context.Context, login string) (entity.User, error)
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
