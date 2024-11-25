package authhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VadimOcLock/gophermart/internal/entity"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"github.com/VadimOcLock/gophermart/internal/usecase/authusecase"
)

type AuthHandler struct {
	AuthUseCase AuthUseCase
}

var _ AuthUseCase = (*authusecase.AuthUseCase)(nil)

func New(authUseCase AuthUseCase) AuthHandler {
	return AuthHandler{
		AuthUseCase: authUseCase,
	}
}

func (h AuthHandler) Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

		return
	}

	var dto entity.UserDTO
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(res, errorz.ErrMsgInvalidRequestFormat, http.StatusBadRequest)

		return
	}
	defer req.Body.Close()

	token, err := h.AuthUseCase.Register(req.Context(), dto)
	if err != nil {
		if errors.Is(err, errorz.ErrLoginAlreadyTaken) {
			http.Error(res, err.Error(), http.StatusConflict)

			return
		}
		if errors.Is(err, errorz.ErrLoginPasswordValidate) {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}
		http.Error(res, "internal server error", http.StatusInternalServerError)

		return
	}

	res.Header().Set("Authorization", "Bearer "+token)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("User successfully registered and authenticated."))
}

func (h AuthHandler) Login(res http.ResponseWriter, req *http.Request) {

}
