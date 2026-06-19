package domain

import "errors"

var (
	ErrEmptyTitle            = errors.New("title is empty")
	ErrNotUniqueTitle        = errors.New("title is not unique")
	ErrEmptyDecription       = errors.New("description is empty")
	ErrInvalidDeadlineFormat = errors.New("invalid deadline format")
	ErrDeadlineInPast        = errors.New("deadline cannot be in the past")
	ErrInvalidStatus         = errors.New("invalid status") // недопустимый статус
	ErrCreatorNotExists      = errors.New("creator does not exists")
	ErrTaskNotFound          = errors.New("task not found")
	ErrToMuchArgs            = errors.New("need only one arg")
	ErrNotEnoughArgs         = errors.New("not enough args")
	ErrInternal              = errors.New("internal error")
)
