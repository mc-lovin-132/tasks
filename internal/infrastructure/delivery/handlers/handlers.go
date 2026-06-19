package handlers

import (
	"context"
	"time"

	"github.com/mc-lovin-132/tasks/internal/domain"
	"github.com/mc-lovin-132/tasks/pb"
)

type service interface {
	Create(ctx context.Context, data *domain.Task) (int, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id *int, title *string) (*domain.Task, error)
	List(ctx context.Context, status *string, isDeadline *bool, creatorID *int) ([]*domain.Task, error)
	Update(ctx context.Context, id int, title, description *string, deadline *time.Time, statusID *int) (int, error)
}

type Handler struct {
	pb.UnimplementedTaskServiceServer
	service service
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	data, err := createReqToDomain(in)
	if err != nil {
		return nil, errorMapper(err)
	}
	id, err := h.service.Create(ctx, data)
	if err != nil {
		return nil, errorMapper(err)
	}
	return createResFromDomain(id), nil
}

func (h *Handler) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	var creatorID *int
	if in.CreatorId != nil {
		cID := int(*in.CreatorId)
		creatorID = &cID
	}
	data, err := h.service.List(ctx, in.Status, in.Deadline, creatorID)
	if err != nil {
		return nil, errorMapper(err)
	}
	return listResFromDomain(data), nil
}

func (h *Handler) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var id *int
	var title *string
	switch in.Selector.(type) {
	case *pb.GetRequest_Id:
		i := int(in.GetId())
		id = &i
	case *pb.GetRequest_Title:
		s := in.GetTitle()
		title = &s
	}
	data, err := h.service.Get(ctx, id, title)
	if err != nil {
		return nil, errorMapper(err)
	}
	return getResFromDomain(data), nil
}
func (h *Handler) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var dd *time.Time
	var statusID *int

	if in.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *in.Deadline)
		if err != nil {
			return nil, domain.ErrInvalidDeadlineFormat
		}
		dd = &t
	}
	if in.StatusId != nil {
		sid := int(*in.StatusId)
		statusID = &sid
	}

	id, err := h.service.Update(ctx, int(in.Id), in.Title, in.Description, dd, statusID)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.UpdateResponse{Id: int64(id)}, nil
}
func (h *Handler) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := h.service.Delete(ctx, int(in.Id))
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.DeleteResponse{}, nil
}
