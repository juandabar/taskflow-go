package handler

import (
	"encoding/json"
	"net/http"

	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/dto"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/httputil"
	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/usecase/auth"
)

type AuthHandler struct {
	registerUseCase *auth.RegisterUserUseCase
	loginUseCase    *auth.LoginUserUseCase
}

func NewAuthHandler(register *auth.RegisterUserUseCase, login *auth.LoginUserUseCase) *AuthHandler {
	return &AuthHandler{
		registerUseCase: register,
		loginUseCase:    login,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, apperror.NewValidationError("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		httputil.WriteError(w, apperror.NewValidationError(err.Error()))
		return
	}

	output, err := h.registerUseCase.Execute(r.Context(), auth.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, dto.RegisterResponse{
		ID:        output.User.ID,
		Name:      output.User.Name,
		Email:     output.User.Email,
		Role:      string(output.User.Role),
		CreatedAt: output.User.CreatedAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, apperror.NewValidationError("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		httputil.WriteError(w, apperror.NewValidationError(err.Error()))
		return
	}

	output, err := h.loginUseCase.Execute(r.Context(), auth.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.LoginResponse{
		Token: output.Token,
	})
}
