package entity

import "time"

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusArchived ProjectStatus = "ARCHIVED"
)

type Project struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
	Status      ProjectStatus
	CreatedAt   time.Time
}
