package entity

import (
	"time"

	"github.com/juandabar/taskflow-go/internal/domain/valueobject"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         valueobject.UserRole
	CreatedAt    time.Time
}
