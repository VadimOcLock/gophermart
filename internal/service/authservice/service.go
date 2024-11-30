package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserStore UserStore
}

func NewAuthService(userStore UserStore) AuthService {
	return AuthService{UserStore: userStore}
}

func (s AuthService) CreateSession(ctx context.Context, userID uint64, token string, expiresAt time.Time) (uint64, error) {
	return s.UserStore.CreateSession(ctx, CreateSessionParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
}

func (s AuthService) IsLoginAvailable(ctx context.Context, login string) (bool, error) {
	exists, err := s.UserStore.UserExistsByLogin(ctx, login)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func (s AuthService) CreateUser(ctx context.Context, login string, password string) (uint64, error) {
	passHash, err := s.hashPassword(password)
	if err != nil {
		return 0, err
	}
	id, err := s.UserStore.CreateUser(ctx, CreateUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s AuthService) FindUserByLogin(ctx context.Context, login string) (entity.User, error) {
	user, err := s.UserStore.FindUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errorz.ErrInvalidLoginPasswordPair
		}

		return entity.User{}, err
	}

	return user, nil
}
