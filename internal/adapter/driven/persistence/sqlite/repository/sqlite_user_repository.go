package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Save(ctx context.Context, user entity.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("saving user: %w", err)
	}
	return nil
}

func (r *SQLiteUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, role, created_at 
		FROM users WHERE email = ?`,
		email,
	)
	return scanUser(row)
}

func (r *SQLiteUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, role, created_at
		FROM users WHERE id = ?`,
		id,
	)
	return scanUser(row)
}

func (r *SQLiteUserRepository) List(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, email, password_hash, role, created_at
		FROM users`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*entity.User, error) {
	var user entity.User
	var role string

	err := scanner.Scan(
		&user.ID, &user.Name, &user.Email,
		&user.PasswordHash, &role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NewNotFoundError("user", "")
	}
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	user.Role = valueobject.UserRole(role)
	return &user, nil
}
