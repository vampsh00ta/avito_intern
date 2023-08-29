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
	AddSegmentsToUser(ctx context.Context, userId int, segments ...*Segment) error
	AddSlugIdToUsers(ctx context.Context, segment Segment, ids ...int) error

	GetUsersSegments(ctx context.Context, userId int) (*[]Segment, error)
	DeleteSegmentsFromUser(ctx context.Context, userId int, segment ...*Segment) error
	GetUserIds(ctx context.Context) ([]int, error)
}
type SegmentRepository interface {
	CreateSegment(ctx context.Context, segment Segment) (int, error)
	DeleteSegment(ctx context.Context, segment Segment) error
}
type History interface {
	AddToHistoryUserSlugs(ctx context.Context, tx pgx.Tx, userId int, operationType bool, segments ...*Segment) error
	AddToHistorySlugUsers(ctx context.Context, tx pgx.Tx, segment Segment, operationType bool, ids ...int) error
	GetHistoryById(ctx context.Context, userId int, year, month int) (*[]HistoryRow, error)
	GetHistoryAll(ctx context.Context, year, month int) (*[]HistoryRow, error)
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
