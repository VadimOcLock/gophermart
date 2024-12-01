package balancehandler

import (
	"encoding/json"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/usecase/balanceusecase"
	"net/http"
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
		http.Error(res, "Unauthorized", http.StatusUnauthorized)

		return
	}
	balance, err := h.BalanceUseCase.FindBalance(req.Context(), userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}
	response, err := json.Marshal(balance)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
