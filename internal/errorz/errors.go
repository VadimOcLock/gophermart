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
	ErrInvalidOrderNumber                = errors.New(`invalid order number`)
	ErrNotEnoughFundsOnBalance           = errors.New("there are not enough funds on the balance")
	ErrUnauthorized                      = errors.New(`unauthorized`)
	ErrExpiredToken                      = errors.New(`expired token`)
)

const (
	ErrMsgOnlyPOSTMethodAcceptMsg = "Only POST method accept"
	ErrMsgInvalidRequestFormatMsg = "Invalid request format"
	ErrInternalServerErrorMsg     = "Internal server error"
	ErrUnauthorizedMsg            = "Unauthorized"
	ErrNoDataToResponseMsg        = "No data to response"
)
