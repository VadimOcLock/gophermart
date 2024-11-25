package authhandler

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/entity"
)

type AuthUseCase interface {
	Register(ctx context.Context, dto entity.UserDTO) (string, error)
}
