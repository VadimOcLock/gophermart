package errorz

import "errors"

var (
	ErrInvalidAddressFormat  = errors.New("invalid address format, expected host:port")
	ErrLoginAlreadyTaken     = errors.New("login already taken")
	ErrLoginPasswordValidate = errors.New("invalid login password format")
)

const (
	ErrMsgOnlyPOSTMethodAccept = "only POST method accept"
	ErrMsgInvalidRequestFormat = "invalid request format"
)
