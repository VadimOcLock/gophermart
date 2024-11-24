package errorz

import "errors"

var (
	ErrInvalidAddressFormat = errors.New("invalid address format, expected host:port")
)
