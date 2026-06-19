package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mc-lovin-132/tasks/internal/domain"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

// неправильный статус
func (r *repository) Create(ctx context.Context, data *domain.Task) (int, error) {
	var id int
	model := fromDomain(data)
	result, err := r.db.NamedQueryContext(ctx, createQuery(), model)
	if err != nil {
		return 0, errorMapper(err)
	}
	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return 0, errorMapper(err)
		}
	} else {
		return 0, errorMapper(sql.ErrNoRows)
	}
	return id, nil
}

func (r *repository) List(ctx context.Context, status *string, isDeadline *bool, creatorID *int) ([]*domain.Task, error) {
	query, args, err := listQuery(status, isDeadline, creatorID)
	if err != nil {
		return nil, errorMapper(err)
	}
	var taskList []*TaskModel
	err = r.db.SelectContext(ctx, &taskList, query, args...)
	if err != nil {
		return nil, errorMapper(err)
	}
	return listToDomain(taskList), nil
}

// нот фаунд
func (r *repository) Get(ctx context.Context, id *int, title *string) (*domain.Task, error) {
	query, args, err := getQuery(title, id)
	fmt.Println(query, args)
	if err != nil {
		return nil, errorMapper(err)
	}
	var task TaskModel
	err = r.db.GetContext(ctx, &task, query, args...)
	if err != nil {
		return nil, errorMapper(err)
	}

	return toDomain(&task), nil
}

// статус не найден
// таска не найдена
func (r *repository) Update(ctx context.Context, id int, title, description *string, deadline *time.Time, statusID *int) (int, error) {
	query, args, err := updateQuery(id, title, description, deadline, statusID)
	if err != nil {
		return 0, errorMapper(err)
	}
	rows, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errorMapper(err)
	}
	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 {
		return 0, errorMapper(sql.ErrNoRows)
	}
	return id, nil
}

// таска не найдена
func (r *repository) Delete(ctx context.Context, id int) error {
	rows, err := r.db.ExecContext(ctx, deleteQuery(), id)
	if err != nil {
		return errorMapper(sql.ErrNoRows)
	}
	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 {
		return errorMapper(sql.ErrNoRows)
	}
	return nil
}
