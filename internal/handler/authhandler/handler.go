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

// Register godoc
// @Summary		Регистрация нового пользователя
// @Description	Регистрация пользователя по логину и паролю. Логины должны быть уникальными. После успешной регистрации происходит автоматическая аутентификация.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		entity.UserDTO	true	"Данные пользователя для регистрации"
// @Success		200		{string}	string			"Пользователь успешно зарегистрирован и аутентифицирован."
// @Failure		400		{string}	string			"Неверный формат запроса."
// @Failure		409		{string}	string			"Логин уже занят."
// @Failure		500		{string}	string			"Внутренняя ошибка сервера."
// @Router			/api/user/register [post]
func (h AuthHandler) Register(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAcceptMsg, http.StatusMethodNotAllowed)

		return
	}

	var dto entity.UserDTO
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(res, errorz.ErrMsgInvalidRequestFormatMsg, http.StatusBadRequest)

		return
	}
	defer req.Body.Close()

	token, err := h.AuthUseCase.Register(req.Context(), dto)
	switch {
	case errors.Is(err, errorz.ErrLoginAlreadyTaken):
		http.Error(res, err.Error(), http.StatusConflict)

		return
	case errors.Is(err, errorz.ErrLoginPasswordValidate):
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	case err != nil:
		http.Error(res, errorz.ErrInternalServerErrorMsg, http.StatusInternalServerError)

		return
	}

	cookie := HTTPCookie(token)
	http.SetCookie(res, cookie)

	res.Header().Set("Authorization", "Bearer "+token)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("User successfully registered and authenticated."))
}

// Login godoc
// @Summary		Аутентификация пользователя
// @Description	Аутентификация пользователя по логину и паролю.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		entity.UserDTO	true	"Данные пользователя для входа"
// @Success		200		{string}	string			"Пользователь успешно аутентифицирован."
// @Failure		400		{string}	string			"Неверный формат запроса."
// @Failure		401		{string}	string			"Неверная пара логин/пароль."
// @Failure		500		{string}	string			"Внутренняя ошибка сервера."
// @Router			/api/user/login [post]
func (h AuthHandler) Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAcceptMsg, http.StatusMethodNotAllowed)

		return
	}
	var dto entity.UserDTO
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(res, errorz.ErrMsgInvalidRequestFormatMsg, http.StatusBadRequest)

		return
	}
	defer req.Body.Close()
	token, err := h.AuthUseCase.Login(req.Context(), dto)
	switch {
	case errors.Is(err, errorz.ErrLoginPasswordValidate):
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	case errors.Is(err, errorz.ErrInvalidLoginPasswordPair):
		http.Error(res, err.Error(), http.StatusUnauthorized)

		return
	case err != nil:
		http.Error(res, errorz.ErrInternalServerErrorMsg, http.StatusInternalServerError)

		return
	}

	cookie := HTTPCookie(token)
	http.SetCookie(res, cookie)

	res.Header().Set("Authorization", "Bearer "+token)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("User successfully authenticated and logged in."))
}

func HTTPCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
