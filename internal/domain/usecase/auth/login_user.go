package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	User  entity.User
	Token string
}

type LoginUserUseCase struct {
	userRepo  userRepository
	jwtSecret string
}

func NewLoginUserUseCase(userRepo userRepository, jwtSecret string) *LoginUserUseCase {
	return &LoginUserUseCase{userRepo: userRepo, jwtSecret: jwtSecret}
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

	token, err := generateToken(user.ID, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	return &LoginUserOutput{User: *user, Token: token}, nil
}

func generateToken(userID, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
