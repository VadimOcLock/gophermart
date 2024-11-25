package entity

import (
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/pkg/validation"
	"time"
)

type UserDTO struct {
	Login    string `json:"login" validate:"required,min=6,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
}

func (dto *UserDTO) Validate() error {
	if err := validation.GetInstance().ValidateStruct(dto); err != nil {
		return errorz.ErrLoginPasswordValidate
	}

	return nil
}

type User struct {
	ID           uint64    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
