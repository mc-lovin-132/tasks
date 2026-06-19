package repository

import (
	"time"

	"github.com/mc-lovin-132/tasks/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

// SELECT
// t.id
// t.title
// t.description
// t.deadline
// t.creator_id
// s.id AS status_id
// s.name AS status_name
// FROM tasks AS t
// LEFT JOIN statuses AS s ON t.status = s.id
// WHERE s.name = $1 AND t.deadline < CURRENT_TIMESTAMP AND creator_id = $2
func listQuery(status *string, isDeadline *bool, creatorID *int) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"t.id", "t.title", "t.description", "t.deadline",
		"t.creator_id", "s.id AS status_id", "s.name AS status_name",
	).
		From("tasks AS t").
		LeftJoin("statuses AS s ON t.status = s.id")

	if status != nil {
		query = query.Where(sq.Eq{"s.name": *status})
	}
	if isDeadline != nil {
		if *isDeadline {
			query = query.Where("t.deadline < CURRENT_TIMESTAMP")
		} else {
			query = query.Where("t.deadline >= CURRENT_TIMESTAMP")
		}
	}
	if creatorID != nil {
		query = query.Where(sq.Eq{"t.creator_id": *creatorID})
	}

	return query.ToSql()
}

// SELECT
// t.id
// t.title
// t.description
// t.deadline
// t.creator_id
// s.id AS status_id
// s.name AS status_name
// FROM tasks AS t
// LEFT JOIN statuses AS s ON t.status = s.id
// WHERE t.title = $1 / id = $1
func getQuery(title *string, id *int) (string, []interface{}, error) {
	if title != nil && id != nil {
		return "", nil, domain.ErrToMuchArgs
	} else if title == nil && id == nil {
		return "", nil, domain.ErrNotEnoughArgs
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"t.id", "t.title", "t.description", "t.deadline",
		"t.creator_id", "s.id AS status_id", "s.name AS status_name",
	).
		From("tasks AS t").
		LeftJoin("statuses AS s ON t.status = s.id")

	if title != nil {
		query = query.Where(sq.Eq{"t.title": *title})
	} else {
		query = query.Where(sq.Eq{"t.id": *id})
	}

	return query.ToSql()
}

// INSERT INTO tasks (
// title,
// description,
// deadline,
// status,
// creator_id
// ) VALUES (
// :title,
// :decription,
// :deadline,
// :status_id,
// :creator_id
// ) RETURNING id;
func createQuery() string {
	return `
	INSERT INTO tasks (
		title, 
		description, 
		deadline, 
		status, 
		creator_id
	) VALUES (
		:title, 
		:description, 
		:deadline, 
		:status_id, 
		:creator_id
	) RETURNING id;`
}

// UPDATE tasks SET
// title = $1,
// description = $2,
// deadline = $3,
// status = $4
// WHERE id = $5;
func updateQuery(id int, title, description *string, deadline *time.Time, statusID *int) (string, []interface{}, error) {
	if title == nil && description == nil && deadline == nil && statusID == nil {
		return "", nil, domain.ErrNotEnoughArgs
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("tasks")
	if title != nil {
		builder = builder.Set("title", *title)
	}
	if description != nil {
		builder = builder.Set("description", *description)
	}
	if deadline != nil {
		builder = builder.Set("deadline", *deadline)
	}
	if statusID != nil {
		builder = builder.Set("status", *statusID)
	}
	return builder.Where(sq.Eq{"id": id}).ToSql()
}

// DELETE FROM tasks WHERE id = $1;
func deleteQuery() string { return `DELETE FROM tasks WHERE id = $1;` }
