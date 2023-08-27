package repository

import (
	postgresql "avito/pkg/client"
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	UserRepository
	SegmentRepository
	History
}
type UserRepository interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, id int) error
	AddSegmentsToUser(ctx context.Context, userId int, slugs ...any) error
	GetUsersSegments(ctx context.Context, userId int) ([]Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, userId int, slugs ...any) error
}
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
}
type History interface {
	AddToHistory(ctx context.Context, tx pgx.Tx, userId int, operationType bool, slugs ...any) error
	GetHistory(ctx context.Context, userId int, year, month int) ([]HistoryRow, error)
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
