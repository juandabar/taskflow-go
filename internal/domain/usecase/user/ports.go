package user

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type userRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
}
