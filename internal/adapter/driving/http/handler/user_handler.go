package handler

import (
	"net/http"

	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/dto"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/httputil"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/middleware"
	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/usecase/user"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type UserHandler struct {
	getUserByIdUseCase *user.GetUserByIdUseCase
	getUsersUseCase    *user.GetUsersUseCase
}

func NewUserHandler(getUserById *user.GetUserByIdUseCase, getUsersUseCase *user.GetUsersUseCase) *UserHandler {
	return &UserHandler{
		getUserByIdUseCase: getUserById,
		getUsersUseCase:    getUsersUseCase,
	}
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	output, err := h.getUserByIdUseCase.Execute(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.GetUserResponse{
		ID:        output.User.ID,
		Name:      output.User.Name,
		Email:     output.User.Email,
		Role:      string(output.User.Role),
		CreatedAt: output.User.CreatedAt,
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	roleStr, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		httputil.WriteError(w, apperror.NewUnauthorizedError("missing role in context"))
	}

	role := valueobject.UserRole(roleStr)

	output, err := h.getUsersUseCase.Execute(r.Context(), role)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	responses := make([]dto.GetUserResponse, 0, len(output.Users))
	for _, u := range output.Users {
		responses = append(responses, dto.GetUserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      string(u.Role),
			CreatedAt: u.CreatedAt,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, responses)
}
