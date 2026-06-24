package auth

import (
	"context"
	"testing"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name         string
		input        RegisterUserInput
		wantError    bool
		whantErrType error
	}{
		{
			name: "success",
			input: RegisterUserInput{
				Name:     "Juan",
				Email:    "juan@example.com",
				Password: "securepassword123",
			},
			wantError: false,
		},
		{
			name: "duplicate email",
			input: RegisterUserInput{
				Name:     "Juan",
				Email:    "juan@example.com",
				Password: "securepassword123",
			},
			wantError:    true,
			whantErrType: &apperror.ConflictError{},
		},
	}

	repo := newInMemoryUserRepository()
	uc := NewRegisterUserUseCase(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := uc.Execute(context.Background(), tt.input)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.whantErrType != nil {
					if _, ok := err.(*apperror.ConflictError); !ok {
						t.Fatalf("expected ConflictError, got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output.User.Role != valueobject.RoleMember {
				t.Fatalf("expected role MEMBER, go %s", output.User.Role)
			}
			if tt.input.Email != output.User.Email {
				t.Fatalf("expected email %s, go %s", tt.input.Email, output.User.Email)
			}
		})
	}
}
