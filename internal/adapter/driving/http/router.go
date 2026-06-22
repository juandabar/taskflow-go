package http

import (
	"net/http"

	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/handler"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("GET /users/{id}", middleware.AuthGuard(jwtSecret)(http.HandlerFunc(userHandler.GetUserById)))

	return mux
}
