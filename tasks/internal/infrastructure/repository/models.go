package repository

import (
	"database/sql"
	"fmt"
	"tasks/internal/domain"
	"time"

	"github.com/lib/pq"
)

func errorMapper(err error) error {
	fmt.Println(err.Error())
	if pqErr, ok := err.(*pq.Error); ok {
		// нарушения внешнего ключа
		if pqErr.Code == "23503" {
			return fmt.Errorf("%w: %w", domain.ErrInvalidStatus, err)
		}
		// нарушение уникальности
		if pqErr.Code == "23505" {
			return fmt.Errorf("%w: %w", domain.ErrNotUniqueTitle, err)
		}
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: %w", domain.ErrTaskNotFound, err)
	}
	return fmt.Errorf("%w: %w", domain.ErrInternal, err)
}

type TaskModel struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Deadline    time.Time `db:"deadline"`
	CreatorID   int       `db:"creator_id"`
	StatusID    int       `db:"status_id"`
	StatusName  string    `db:"status_name"`
}

// domain -> model
func fromDomain(task *domain.Task) *TaskModel {
	return &TaskModel{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		StatusID:    task.Status.ID,
		CreatorID:   task.CreatorID,
	}
}

// model -> domain
func toDomain(task *TaskModel) *domain.Task {
	return &domain.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status: domain.Status{
			ID:   task.StatusID,
			Name: task.StatusName,
		},
		CreatorID: task.CreatorID,
	}
}

// models -> domains
func listToDomain(tasks []*TaskModel) []*domain.Task {
	lst := make([]*domain.Task, len(tasks))
	for i, task := range tasks {
		lst[i] = toDomain(task)
	}
	return lst
}
