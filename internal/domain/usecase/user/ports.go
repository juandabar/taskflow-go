package user

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type userFinder interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
}

type usersLister interface {
	List(ctx context.Context) ([]entity.User, error)
}
