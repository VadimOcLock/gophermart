package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

var instance *Validator
var once sync.Once

func GetInstance() *Validator {
	once.Do(func() {
		instance = &Validator{
			v: validator.New(),
		}
	})

	return instance
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.v.Struct(s)
}
