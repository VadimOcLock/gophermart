package balancehandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/usecase/balanceusecase"
)

type BalanceHandler struct {
	BalanceUseCase BalanceUseCase
}

var _ BalanceUseCase = (*balanceusecase.BalanceUseCase)(nil)

func New(balanceUseCase BalanceUseCase) *BalanceHandler {
	return &BalanceHandler{
		BalanceUseCase: balanceUseCase,
	}
}

func (h BalanceHandler) GetBalance(res http.ResponseWriter, req *http.Request) {
	if http.MethodGet != req.Method {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, errorz.ErrUnauthorized, http.StatusUnauthorized)

		return
	}
	balance, err := h.BalanceUseCase.FindBalance(req.Context(), userID)
	if err != nil {
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)

		return
	}
	response, err := json.Marshal(balance)
	if err != nil {
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)

		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(response)
}

func (h BalanceHandler) WithdrawBalance(res http.ResponseWriter, req *http.Request) {
	if http.MethodPost != req.Method {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, errorz.ErrUnauthorized, http.StatusUnauthorized)

		return
	}
	var dto entity.Withdraw
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		http.Error(res, errorz.ErrMsgInvalidRequestFormat, http.StatusBadRequest)

		return
	}
	err = h.BalanceUseCase.Withdrawal(req.Context(), userID, dto.Sum, dto.OrderNumber)
	switch {
	case err == nil:
		res.WriteHeader(http.StatusOK)
	case errors.Is(err, errorz.ErrInvalidOrderNumber):
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, errorz.ErrNotEnoughFundsOnBalance):
		http.Error(res, err.Error(), http.StatusPaymentRequired)
	default:
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)
	}
}

func (h BalanceHandler) GetWithdrawals(res http.ResponseWriter, req *http.Request) {
	if http.MethodGet != req.Method {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, errorz.ErrUnauthorized, http.StatusUnauthorized)

		return
	}
	response, err := h.BalanceUseCase.FindWithdrawals(req.Context(), userID)
	if err != nil {
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)
	}
	switch {
	case err == nil:
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(response)
	case errors.Is(err, errorz.ErrUserHasNoWithdrawals):
		http.Error(res, err.Error(), http.StatusNoContent)
	default:
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)
	}
}
