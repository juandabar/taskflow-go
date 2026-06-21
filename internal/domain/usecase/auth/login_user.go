package auth

import (
	"context"
	"fmt"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	User entity.User
}

type LoginUserUseCase struct {
	userRepo userRepository
}

func NewLoginUserUseCase(userRepo userRepository) *LoginUserUseCase {
	return &LoginUserUseCase{userRepo: userRepo}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if _, ok := err.(*apperror.NotFoundError); ok {
			return nil, apperror.NewValidationError("invalid credentials")
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, apperror.NewValidationError("invalid credentials")
	}

	return &LoginUserOutput{User: *user}, nil
}
