package auth

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user entity.User) error
}
