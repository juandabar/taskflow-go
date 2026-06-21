package repository

import (
	"context"

	"github.com/juandabar/taskflow-go/internal/domain/entity"
)

type CommentRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Comment, error)
	Save(ctx context.Context, comment entity.Comment) error
	Delete(ctx context.Context, id string) error
	ListByTask(ctx context.Context, taskID string) ([]entity.Comment, error)
}
