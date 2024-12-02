package middleware

import (
	"context"
	"github.com/VadimOcLock/gophermart/internal/errorz"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey  = "user_id"
	CookieName = "jwt_token"
)

type JWTValidator struct {
	JWTSecretKey string
}

func NewJWTValidator(jwtSecretKey string) *JWTValidator {
	return &JWTValidator{JWTSecretKey: jwtSecretKey}
}

func (v *JWTValidator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(v.JWTSecretKey), nil
	})
}

func (v *JWTValidator) ExtractAndVerifyToken(r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
			tokenString := bearerToken[1]
			return v.verifyToken(tokenString, r)
		}
	}

	cookie, err := r.Cookie(CookieName)
	if err == nil {
		tokenString := cookie.Value
		return v.verifyToken(tokenString, r)
	}

	return nil, errorz.ErrUnauthorized
}

func (v *JWTValidator) verifyToken(tokenString string, r *http.Request) (context.Context, error) {
	token, err := v.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, errorz.ErrUnauthorized
		}

		expiresIn := time.Unix(int64(exp), 0)
		if time.Now().After(expiresIn) {
			return nil, errorz.ErrExpiredToken
		}

		userID, ok := claims[UserIDKey].(float64)
		if !ok {
			return nil, errorz.ErrUnauthorized
		}

		ctx := context.WithValue(r.Context(), UserIDKey, uint64(userID))
		return ctx, nil
	}

	return nil, errorz.ErrUnauthorized
}

func JWTAuthMiddleware(jwtSecretKey string) func(next http.Handler) http.Handler {
	jwtValidator := NewJWTValidator(jwtSecretKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, err := jwtValidator.ExtractAndVerifyToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)

	return userID, ok
}
