package domain

import "time"

type Status struct {
	ID   int
	Name string
}

type Task struct {
	ID          int
	Title       string
	Description string
	Deadline    time.Time
	Status      Status
	CreatorID   int
}
