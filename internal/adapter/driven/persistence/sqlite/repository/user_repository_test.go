package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
	"github.com/juandabar/taskflow-go/internal/infrastructure/database"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := database.NewSQLiteConnection(":memory:")
	if err != nil {
		t.Fatalf("failed to setup test database %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestSQLiteUserRepository_Save_And_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteUserRepository(db)

	user := entity.User{
		ID:           "test-id-001",
		Name:         "Juan",
		Email:        "juan@example.com",
		PasswordHash: "hashedpassword",
		Role:         valueobject.RoleMember,
		CreatedAt:    time.Now().UTC(),
	}

	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, found.ID)
	}
	if found.Email != user.Email {
		t.Errorf("expected email %s, go %s", user.Email, found.Email)
	}
	if found.Role != user.Role {
		t.Errorf("expected role %s, go %s", user.Role, found.Role)
	}

}
