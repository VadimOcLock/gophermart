package authusecase

import (
	"context"
	"time"

	"github.com/VadimOcLock/gophermart/internal/entity"
)

type AuthService interface {
	IsLoginAvailable(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password string) (uint64, error)
	CreateSession(ctx context.Context, userID uint64, token string, expiresAt time.Time) (uint64, error)
	FindUserByLogin(ctx context.Context, login string) (entity.User, error)
}

type CreateUserParams struct {
	Login    string
	Password string
}

type CreateSessionParams struct {
	UserID    uint64
	Token     string
	ExpiresAt time.Time
}
