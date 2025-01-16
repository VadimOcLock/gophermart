package orderhandler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/middleware"
	"github.com/VadimOcLock/gophermart/internal/usecase/orderusecase"
)

type OrderHandler struct {
	OrderUseCase OrderUseCase
}

var _ OrderUseCase = (*orderusecase.OrderUseCase)(nil)

func New(orderUseCase OrderUseCase) OrderHandler {
	return OrderHandler{OrderUseCase: orderUseCase}
}

// UploadOrder godoc
// @Summary Загрузить номер заказа
// @Description Хендлер доступен только аутентифицированным пользователям.
// @Tags orders
// @Security BearerAuth
// @Accept text/plain
// @Produce json
// @Param order_number body string true "Номер заказа"
// @Success 200 {string} string "Номер заказа уже был загружен этим пользователем."
// @Success 202 {string} string "Новый номер заказа принят в обработку."
// @Failure 400 {string} string "Неверный формат запроса."
// @Failure 401 {string} string "Пользователь не аутентифицирован."
// @Failure 409 {string} string "Номер заказа уже был загружен другим пользователем."
// @Failure 422 {string} string "Неверный формат номера заказа."
// @Failure 500 {string} string "Внутренняя ошибка сервера."
// @Router /api/user/orders [post]
func (h OrderHandler) UploadOrder(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, errorz.ErrUnauthorized, http.StatusUnauthorized)

		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		http.Error(res, errorz.ErrMsgInvalidRequestFormat, http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)
	orderNumber := strings.TrimSpace(string(body))
	err = h.OrderUseCase.UploadOrder(req.Context(), userID, orderNumber)
	switch {
	case err == nil:
		res.WriteHeader(http.StatusAccepted)
	case errors.Is(err, errorz.ErrOrderAlreadyUploadedByUser):
		http.Error(res, err.Error(), http.StatusOK)
	case errors.Is(err, errorz.ErrOrderAlreadyUploadedByAnotherUser):
		http.Error(res, err.Error(), http.StatusConflict)
	case errors.Is(err, errorz.ErrInvalidOrderNumberFormat):
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
	default:
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)
	}
}

// GetOrders godoc
// @Summary Получение списка загруженных номеров заказов
// @Description Хендлер доступен только авторизованному пользователю. Номера заказа в отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Order "Список загруженных номеров заказов"
// @Success 204 "Нет данных для ответа"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/orders [get]
func (h OrderHandler) GetOrders(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
	userID, ok := middleware.UserIDFromContext(req.Context())
	if !ok {
		http.Error(res, errorz.ErrUnauthorized, http.StatusUnauthorized)

		return
	}
	response, err := h.OrderUseCase.FindAllOrders(req.Context(), userID)
	switch {
	case err == nil:
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(response)
	case errors.Is(err, errorz.ErrUserHasNoOrders):
		http.Error(res, errorz.ErrNoDataToResponse, http.StatusNoContent)
	default:
		http.Error(res, errorz.ErrInternalServerError, http.StatusInternalServerError)
	}
}
