package http

import (
	"net/http"

	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/handler"
)

func NewRouter(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	return mux
}
