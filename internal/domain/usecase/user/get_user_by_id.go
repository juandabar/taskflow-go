package user

import (
	"context"
	"fmt"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type GetUserByIdUseCase struct {
	userRepo userRepository
}

type GetUserOutput struct {
	User entity.User
}

func NewGetUserByIdUseCase(userRepo userRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepo: userRepo}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, id string) (*GetUserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if _, ok := err.(*apperror.NotFoundError); ok {
			return nil, apperror.NewNotFoundError("user", id)
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	return &GetUserOutput{
		User: *user,
	}, nil
}
