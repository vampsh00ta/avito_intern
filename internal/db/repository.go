package repository

import "context"

type Repository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, username string) error
	AddTagToUser(ctx context.Context, userId int, slugs ...any) error
	GetSegmentsIds(ctx context.Context, slugs ...any) ([]int, error)
}
