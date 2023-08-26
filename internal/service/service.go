package service

import (
	rep "avito/internal/db"
	"avito/internal/transport/model"
	"context"
)

type Service interface {
	User
	Segment
}
type service struct {
	rep rep.Repository
}

func New(r rep.Repository) Service {
	return service{
		r,
	}
}

type User interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, userId int) error
	AddSegmentsToUser(ctx context.Context, data model.RequestAddOrDeleteSegmentsToUser) error
	GetUsersSegments(ctx context.Context, userId int) ([]model.Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, data model.RequestAddOrDeleteSegmentsToUser) (err error)
}

type Segment interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
}
