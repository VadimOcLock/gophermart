package authusecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/service/authservice"
	"github.com/VadimOcLock/gophermart/pkg/jwt"
)

type AuthUseCase struct {
	AuthService AuthService
	JWTConfig   JWTConfig
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

var _ AuthService = (*authservice.AuthService)(nil)

func NewAuthUseCase(authService AuthService, jwtCfg JWTConfig) AuthUseCase {
	return AuthUseCase{
		AuthService: authService,
		JWTConfig:   jwtCfg,
	}
}

func (uc AuthUseCase) Register(ctx context.Context, dto entity.UserDTO) (string, error) {
	if err := dto.Validate(); err != nil {
		return "", err
	}

	available, err := uc.AuthService.IsLoginAvailable(ctx, dto.Login)
	if err != nil {
		return "", err
	}
	if !available {
		return "", errorz.ErrLoginAlreadyTaken
	}

	userID, err := uc.AuthService.CreateUser(ctx, dto.Login, dto.Password)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(uc.JWTConfig.TokenDuration)
	token, err := jwt.Generate(userID, expiresAt, uc.JWTConfig.SecretKey)
	if err != nil {
		return "", err
	}

	if _, err = uc.AuthService.CreateSession(ctx, userID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func (uc AuthUseCase) Login(ctx context.Context, dto entity.UserDTO) (string, error) {
	if err := dto.Validate(); err != nil {
		return "", err
	}
	user, err := uc.AuthService.FindUserByLogin(ctx, dto.Login)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return "", errorz.ErrInvalidLoginPasswordPair
	}
	expiresAt := time.Now().Add(uc.JWTConfig.TokenDuration)
	token, err := jwt.Generate(user.ID, expiresAt, uc.JWTConfig.SecretKey)
	if err != nil {
		return "", err
	}
	if _, err = uc.AuthService.CreateSession(ctx, user.ID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}
