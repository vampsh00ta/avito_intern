package service

import (
	rep "avito/internal/db"
	"avito/internal/ttl"
	"context"
)

type Service interface {
	User
	Segment
	History
}
type service struct {
	rep rep.Repository
	ttl ttl.TTL
}

func New(r rep.Repository, ttl ttl.TTL) Service {
	return service{
		r,
		ttl,
	}
}

type User interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, userId int) error
	AddSegmentsToUser(ctx context.Context, userId int, segments ...*rep.Segment) error
	GetUsersSegments(ctx context.Context, userId int) ([]rep.Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, userId int, segments ...*rep.Segment) (err error)
}

type Segment interface {
	CreateSegment(ctx context.Context, segment rep.Segment) error
	CreateSegmentPercent(ctx context.Context, segment rep.Segment) (*[]rep.User, error)
	DeleteSegment(ctx context.Context, segment rep.Segment) error
}
type History interface {
	GetHistory(ctx context.Context, userId, year, month int) (*[]rep.HistoryRow, error)
}
