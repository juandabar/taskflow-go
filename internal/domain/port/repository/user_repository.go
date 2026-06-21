package repository

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user entity.User) error
	List(ctx context.Context) ([]entity.User, error)
}
