package repository

import (
	"avito/internal/transport/model"
	postgresql "avito/pkg/client"
	"context"
	"go.uber.org/zap"
)

type Repository interface {
	UserRepository
	SegmentRepository
}
type UserRepository interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, id int) error
	AddSegmentsToUser(ctx context.Context, userId int, slugs ...any) error
	GetUsersSegments(ctx context.Context, userId int) ([]model.Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, userId int, slugs ...any) error
}
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	GetSegmentsIds(ctx context.Context, tx interface{}, slugs ...any) ([]int, error)
}
type Db struct {
	client postgresql.Client
	log    *zap.SugaredLogger
}

func New(client postgresql.Client, logger *zap.SugaredLogger) Repository {
	return &Db{
		client,
		logger,
	}
}
