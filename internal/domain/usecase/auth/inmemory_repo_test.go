package auth

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type inMemoryUserRepository struct {
	users map[string]entity.User
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]entity.User),
	}
}

func (r *inMemoryUserRepository) Save(ctx context.Context, user entity.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *inMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, apperror.NewNotFoundError("user", email)
	}
	return &user, nil
}
