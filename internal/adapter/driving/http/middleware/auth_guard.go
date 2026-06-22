package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/httputil"
	"github.com/juandabar/taskflow-go/internal/domain/apperror"
)

type contextKey string

const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

func AuthGuard(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				httputil.WriteError(w, apperror.NewValidationError("missing or invalid authorization header"))
				return
			}

			tokenString := strings.TrimPrefix(header, "Bearer ")

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, apperror.NewValidationError("invalid token signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				httputil.WriteError(w, apperror.NewValidationError("invalid or expired token"))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				httputil.WriteError(w, apperror.NewValidationError("invalid token claims"))
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				httputil.WriteError(w, apperror.NewValidationError("invalid token subject"))
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				httputil.WriteError(w, apperror.NewValidationError("invalid token role"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
