package entity

import "time"

type Comment struct {
	ID        string
	Content   string
	TaskID    string
	AuthorID  string
	CreatedAt time.Time
}
