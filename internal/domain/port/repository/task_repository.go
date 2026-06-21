package repository

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type TaskFilter struct {
	Status     *valueobject.TaskStatus
	Priority   *valueobject.Priority
	AssigneeID *string
}

type TaskRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Task, error)
	Save(ctx context.Context, task entity.Task) error
	Update(ctx context.Context, task entity.Task) error
	ListByProject(ctx context.Context, projectID string, filter TaskFilter) ([]entity.Task, error)
}
