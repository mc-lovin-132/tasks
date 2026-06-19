package app

import (
	"context"
	"time"

	"github.com/mc-lovin-132/tasks/internal/domain"
)

type userServiceMock struct{}

func (u *userServiceMock) IsUserExists(id int) bool { return true }

type eventSendenerMock struct{}

func (e *eventSendenerMock) TaskCreated(data *domain.Task) error {
	if data != nil {
		return nil
	}
	return nil
}
func (e *eventSendenerMock) DeadlineSoon(data *domain.Task, until time.Duration) error {
	if data != nil && until > 0 {
		return nil
	}
	return nil
}

type repoMock struct{}

func (r *repoMock) Create(ctx context.Context, data *domain.Task) (int, error) { return 0, nil }
func (r *repoMock) Delete(ctx context.Context, id int) error                   { return nil }
func (r *repoMock) Get(ctx context.Context, id *int, title *string) (*domain.Task, error) {
	return &domain.Task{}, nil
}
func (r *repoMock) List(ctx context.Context, status *string, isDeadline *bool, creatorID *int) ([]*domain.Task, error) {
	return []*domain.Task{
		{}, {}, {},
	}, nil
}
func (r *repoMock) Update(ctx context.Context, id int, title, description *string, deadline *time.Time, statusID *int) (int, error) {
	return 0, nil
}
