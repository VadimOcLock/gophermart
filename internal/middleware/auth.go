package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

const UserIDKey = "user_id"

func JWTAuthMiddleware(jwtSecretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(res, "authorization header is missing", http.StatusUnauthorized)

				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(res, "invalid authorization header format", http.StatusBadRequest)

				return
			}

			tokenString := bearerToken[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(jwtSecretKey), nil
			})
			if err != nil {
				http.Error(res, "invalid token", http.StatusUnauthorized)

				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, ok := claims[UserIDKey].(uint64)
				if !ok {
					http.Error(res, "missing subject claim", http.StatusUnauthorized)

					return
				}

				exp, ok := claims["exp"].(float64)
				if !ok {
					http.Error(res, "missing expiry claim", http.StatusUnauthorized)

					return
				}

				expiresIn := time.Unix(int64(exp), 0)
				if time.Now().After(expiresIn) {
					http.Error(res, "token expired", http.StatusUnauthorized)

					return
				}

				ctx := context.WithValue(req.Context(), UserIDKey, userID)
				next.ServeHTTP(res, req.WithContext(ctx))
			} else {
				http.Error(res, "unauthorized access", http.StatusUnauthorized)
			}
		})
	}
}

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)

	return userID, ok
}
