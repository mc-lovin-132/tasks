package service

import (
	"context"
	"tasks/internal/domain"
	"time"
)

type repository interface {
	Create(ctx context.Context, data *domain.Task) (int, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id *int, title *string) (*domain.Task, error)
	List(ctx context.Context, status *string, isDeadline *bool, creatorID *int) ([]*domain.Task, error)
	Update(ctx context.Context, id int, title, description *string, deadline *time.Time, statusID *int) (int, error)
}

type userService interface {
	IsUserExists(id int) bool
}

type eventSendener interface {
	TaskCreated(data *domain.Task) error
	DeadlineSoon(data *domain.Task, until time.Duration) error
}

type service struct {
	userService   userService
	repo          repository
	eventSendener eventSendener
}

func New(repo repository, userService userService, eventSendener eventSendener) *service {
	return &service{
		repo:          repo,
		userService:   userService,
		eventSendener: eventSendener,
	}
}
func (s *service) Create(ctx context.Context, data *domain.Task) (int, error) {
	if !s.userService.IsUserExists(data.CreatorID) {
		return 0, domain.ErrCreatorNotExists
	}
	id, err := s.repo.Create(ctx, data)
	data.ID = id
	if err != nil {
		return 0, err
	}
	err = s.eventSendener.TaskCreated(data)
	if err != nil {
		return 0, err
	}
	return id, err
}
func (s *service) List(ctx context.Context, status *string, isDeadline *bool, creatorID *int) ([]*domain.Task, error) {
	return s.repo.List(ctx, status, isDeadline, creatorID)
}
func (s *service) Get(ctx context.Context, id *int, title *string) (*domain.Task, error) {
	return s.repo.Get(ctx, id, title)
}
func (s *service) Update(ctx context.Context, id int, title, description *string, deadline *time.Time, statusID *int) (int, error) {
	return s.repo.Update(ctx, id, title, description, deadline, statusID)
}
func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
