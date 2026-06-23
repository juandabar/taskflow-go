package user

import (
	"context"
	"fmt"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type GetUsersUseCase struct {
	userRepo usersLister
}

type GetUsersOutput struct {
	Users []entity.User
}

func NewGetUsersUseCase(userRepo usersLister) *GetUsersUseCase {
	return &GetUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUsersUseCase) Execute(ctx context.Context, role valueobject.UserRole) (*GetUsersOutput, error) {
	if !role.IsAdmin() {
		return nil, apperror.NewForbiddenError("only admins can list users")
	}
	users, err := uc.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users use case: %w", err)
	}

	return &GetUsersOutput{
		Users: users,
	}, nil
}
