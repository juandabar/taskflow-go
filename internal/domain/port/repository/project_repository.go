package repository

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type ProjectRepository interface {
	FindById(ctx context.Context, id string) (*entity.Project, error)
	Save(ctx context.Context, project entity.Project) error
	Update(ctx context.Context, project entity.Project) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status *entity.ProjectStatus) ([]entity.Project, error)
}
