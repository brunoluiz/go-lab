package model

import "time"

type List struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type Task struct {
	ID          string
	Title       string
	IsCompleted bool
	CreatedAt   time.Time
	ListID      string
}
