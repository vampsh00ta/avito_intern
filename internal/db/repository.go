package repository

import "context"

type Repository interface {
	UserRepository
	SegmentRepository
}
type UserRepository interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, username string) error
	AddTagsToUser(ctx context.Context, userId int, slugs ...any) error
	GetUsersTags(ctx context.Context, userId int) ([]Segment, error)
}
type SegmentRepository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	AddTagsToUser(ctx context.Context, userId int, slugs ...any) error
	GetSegmentsIds(ctx context.Context, tx interface{}, slugs ...any) ([]int, error)
}
