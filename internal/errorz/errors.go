package errorz

import "errors"

var (
	ErrInvalidAddressFormat              = errors.New("invalid address format, expected host:port")
	ErrLoginAlreadyTaken                 = errors.New("login already taken")
	ErrLoginPasswordValidate             = errors.New("invalid login password format")
	ErrInvalidLoginPasswordPair          = errors.New("invalid  login password pair")
	ErrOrderAlreadyUploadedByUser        = errors.New("order already uploaded by user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrInvalidOrderNumberFormat          = errors.New("invalid order number format")
	ErrUserHasNoOrders                   = errors.New("user has no orders")
	ErrUserHasNoWithdrawals              = errors.New("user has no withdrawals")
)

const (
	ErrMsgOnlyPOSTMethodAccept = "Only POST method accept"
	ErrMsgInvalidRequestFormat = "Invalid request format"
	ErrInternalServerError     = "Internal server error"
	ErrUnauthorized            = "Unauthorized"
	ErrNoDataToResponse        = "No data to response"
)
