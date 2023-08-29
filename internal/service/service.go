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
	AddSegmentsToUser(ctx context.Context, userId int, segments ...*Segment_AddSegmentsToUser) error
	GetUsersSegments(ctx context.Context, userId int) ([]rep.Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, userId int, segments ...*Segment_DeleteSegmentsFromUser) (err error)
}

type Segment interface {
	CreateSegment(ctx context.Context, segment Segment_CreateSegment) error
	CreateSegmentPercent(ctx context.Context, segment Segment_CreateSegment) (*[]User_CreateSegmentPercent, error)
	DeleteSegment(ctx context.Context, segment Segment_DeleteSegment) error
}
type History interface {
	GetHistory(ctx context.Context, userId, year, month int) (*[]rep.HistoryRow, error)
}
