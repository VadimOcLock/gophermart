package entity

import (
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/pkg/validation"
)

type UserDTO struct {
	Login    string `json:"login" validate:"required,min=6,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=8,regexp=^[^ ]+$"`
}

func (dto *UserDTO) Validate() error {
	if err := validation.GetInstance().ValidateStruct(dto); err != nil {
		return errorz.ErrLoginPasswordValidate
	}

	//validate := validator.New()
	//if err := validate.Struct(dto); err != nil {
	//	return errorz.ErrLoginPasswordValidate
	//}

	return nil
}
