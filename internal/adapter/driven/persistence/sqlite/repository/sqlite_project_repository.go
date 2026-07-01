package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/juandabar/taskflow-go/internal/domain/apperror"
	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type SQLiteProjectRepository struct {
	db *sql.DB
}

func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

func (r *SQLiteProjectRepository) Save(ctx context.Context, project entity.Project) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (id, name, description, owner_id, status, created_at)
		VALUES (?, ?)`,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
		project.Status,
		project.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("saving project: %w", err)
	}
	return nil
}

func (r *SQLiteProjectRepository) FindByID(ctx context.Context, id string) (*entity.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, owner_id, status, created_at
		FROM projects
		WHERE id = ?`,
		id,
	)
	return scanProject(row)
}

func scanProject(scanner interface {
	Scan(dest ...any) error
}) (*entity.Project, error) {
	var project entity.Project
	var status string
	var createdAt string

	err := scanner.Scan(
		&project.ID, &project.Name, &project.Description,
		&project.OwnerID, &status, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NewNotFoundError("project", "")
	}
	if err != nil {
		return nil, fmt.Errorf("scanning project: %w", err)
	}

	project.Status = entity.ProjectStatus(status)
	project.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999 +0000 UTC", createdAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}

	return &project, nil
}
