package handlers

import (
	"errors"
	"time"

	"github.com/mc-lovin-132/tasks/internal/domain"
	"github.com/mc-lovin-132/tasks/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorMapper(err error) error {
	if errors.Is(err, domain.ErrEmptyTitle) ||
		errors.Is(err, domain.ErrEmptyDecription) ||
		errors.Is(err, domain.ErrInvalidDeadlineFormat) ||
		errors.Is(err, domain.ErrDeadlineInPast) ||
		errors.Is(err, domain.ErrInvalidStatus) ||
		errors.Is(err, domain.ErrToMuchArgs) ||
		errors.Is(err, domain.ErrNotEnoughArgs) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, domain.ErrNotUniqueTitle) {
		return status.Error(codes.AlreadyExists, err.Error())
	} else if errors.Is(err, domain.ErrCreatorNotExists) || errors.Is(err, domain.ErrTaskNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, domain.ErrInternal) {
		return status.Error(codes.Internal, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}

func statusToDomain(status *pb.Status) domain.Status {
	return domain.Status{
		ID:   int(status.Id),
		Name: status.Name,
	}
}
func statusFromDomain(status domain.Status) *pb.Status {
	return &pb.Status{Id: int64(status.ID), Name: status.Name}
}

// TODO: добавить ограничения длины строковых полей
// title > 0 && < 255
// description > 0 && < 1000
func createReqToDomain(req *pb.CreateRequest) (*domain.Task, error) {
	// 2006-01-02T15:04:05Z
	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		return nil, domain.ErrInvalidDeadlineFormat
	}
	task := &domain.Task{}

	if req.Title == "" {
		return nil, domain.ErrEmptyTitle
	}
	task.Title = req.Title

	if req.Description == "" {
		return nil, domain.ErrEmptyDecription
	}
	task.Description = req.Description

	if time.Now().After(deadline) {
		return nil, domain.ErrDeadlineInPast
	}
	task.Deadline = deadline

	task.CreatorID = int(req.CreatorID)
	task.Status = domain.Status{ID: int(req.StatusID)}
	return task, nil
}
func createResFromDomain(id int) *pb.CreateResponse {
	return &pb.CreateResponse{Id: int64(id)}
}
func domainToPb(data *domain.Task) *pb.Task {
	return &pb.Task{
		Id:          int64(data.ID),
		Title:       data.Title,
		Description: data.Description,
		Deadline:    data.Deadline.Format(time.RFC3339),
		Status:      statusFromDomain(data.Status),
		CreatorID:   int64(data.CreatorID),
	}
}

func listResFromDomain(data []*domain.Task) *pb.ListResponse {
	lst := make([]*pb.Task, len(data))
	for i, t := range data {
		lst[i] = domainToPb(t)
	}
	return &pb.ListResponse{Tasks: lst}
}

func getResFromDomain(data *domain.Task) *pb.GetResponse {
	return &pb.GetResponse{Task: domainToPb(data)}
}
