package orderhandler

import (
	"errors"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/usecase/orderusecase"
	"io"
	"net/http"
	"strings"
)

type OrderHandler struct {
	OrderUseCase OrderUseCase
}

var _ OrderUseCase = (*orderusecase.OrderUseCase)(nil)

func New(orderUseCase OrderUseCase) OrderHandler {
	return OrderHandler{OrderUseCase: orderUseCase}
}

func (h OrderHandler) UploadOrder(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)

		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		http.Error(res, "Invalid request format", http.StatusBadRequest)

		return
	}
	defer req.Body.Close()
	orderNumber := strings.TrimSpace(string(body))
	err = h.OrderUseCase.UploadOrder(req.Context(), userID, orderNumber)
	switch {
	case err == nil:
		res.WriteHeader(http.StatusAccepted)
	case errors.Is(err, errorz.ErrOrderAlreadyUploadedByUser):
		http.Error(res, "User already uploaded by user", http.StatusOK)
	case errors.Is(err, errorz.ErrOrderAlreadyUploadedByAnotherUser):
		http.Error(res, "User already uploaded by another user", http.StatusConflict)
	case errors.Is(err, errorz.ErrInvalidOrderNumberFormat):
		http.Error(res, "Invalid order number format", http.StatusUnprocessableEntity)
	default:
		http.Error(res, "Internal server error", http.StatusInternalServerError)
	}
}
