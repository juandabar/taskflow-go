package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"github.com/juandabar/taskflow-go/internal/domain/port/repository"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	User entity.User
}

type RegisterUserUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUserUseCase(userRepo repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

func (uc RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if _, ok := err.(*apperror.NotFoundError); !ok {
			return nil, err
		}
	}
	if existing != nil {
		return nil, apperror.NewConflictError("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := entity.User{
		ID:           generateID(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         valueobject.RoleMember,
		CreatedAt:    time.Now().UTC(),
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterUserOutput{User: user}, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
