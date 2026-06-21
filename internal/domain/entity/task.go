package entity

import (
	"time"

	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type Task struct {
	ID          string
	Title       string
	Description string
	ProjectID   string
	AssigneeID  string
	Status      valueobject.TaskStatus
	Priority    valueobject.Priority
	DueDate     *time.Time
	CreatedAt   time.Time
}
