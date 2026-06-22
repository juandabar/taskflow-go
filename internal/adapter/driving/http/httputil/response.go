package httputil

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
)

type problemDetail struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func WriteError(w http.ResponseWriter, err error) {
	var pd problemDetail

	switch e := err.(type) {
	case *apperror.NotFoundError:
		pd = problemDetail{
			Type:   "https://taskflow.api/errors/not-found",
			Title:  "Resource not found",
			Status: http.StatusNotFound,
			Detail: e.Error(),
		}
	case *apperror.ValidationError:
		pd = problemDetail{
			Type:   "https://taskflow.api/errors/validation",
			Title:  "Validation error",
			Status: http.StatusBadRequest,
			Detail: e.Error(),
		}
	case *apperror.ConflictError:
		pd = problemDetail{
			Type:   "https://taskflow.api/errors/conflict",
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: e.Error(),
		}
	case *apperror.ForbiddenError:
		pd = problemDetail{
			Type:   "https://taskflow.api/errors/forbidden",
			Title:  "Forbidden",
			Status: http.StatusForbidden,
			Detail: e.Error(),
		}
	default:
		log.Printf("internal error: %v", err)
		pd = problemDetail{
			Type:   "https://taskflow.api/errors/internal",
			Title:  "Internal server error",
			Status: http.StatusInternalServerError,
			Detail: "an unexpected error occurred",
		}
	}

	WriteJSON(w, pd.Status, pd)
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
