package handler

import (
	"net/http"

	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/dto"
	"github.com/juandabar/taskflow-go/internal/adapter/driving/http/httputil"
	"github.com/juandabar/taskflow-go/internal/domain/usecase/user"
)

type UserHandler struct {
	getUserByIdUseCase *user.GetUserByIdUseCase
}

func NewUserHandler(getUserById *user.GetUserByIdUseCase) *UserHandler {
	return &UserHandler{
		getUserByIdUseCase: getUserById,
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
