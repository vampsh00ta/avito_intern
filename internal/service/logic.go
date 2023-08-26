package service

import (
	rep "avito/internal/db"
	"avito/internal/transport/model"
	"context"
)

func (s service) CreateUser(ctx context.Context, username string) error {

	return s.rep.CreateUser(ctx, username)
}

func (s service) DeleteUser(ctx context.Context, userId int) error {

	return s.rep.DeleteUser(ctx, userId)
}

func (s service) AddSegmentsToUser(ctx context.Context, data model.RequestAddOrDeleteSegmentsToUser) error {
	var slugs []any
	for _, slug := range data.Segments {
		slugs = append(slugs, slug.Slug)
	}
	return s.rep.AddSegmentsToUser(ctx, data.Id, slugs...)
}

func (s service) GetUsersSegments(ctx context.Context, userId int) ([]rep.Segment, error) {
	return s.rep.GetUsersSegments(ctx, userId)
}
func (s service) DeleteSegmentsFromUser(ctx context.Context, data model.RequestAddOrDeleteSegmentsToUser) (err error) {
	var slugs []any
	for _, slug := range data.Segments {
		slugs = append(slugs, slug.Slug)
	}
	return s.rep.DeleteSegmentsFromUser(ctx, data.Id, slugs...)
}

func (s service) CreateSegment(ctx context.Context, slug string) error {
	return s.rep.CreateSegment(ctx, slug)

}
func (s service) DeleteSegment(ctx context.Context, slug string) error {
	return s.rep.DeleteSegment(ctx, slug)

}

func (s service) GetHistory(ctx context.Context, userId, year, month int) ([]rep.HistoryRow, error) {
	history, err := s.rep.GetHistory(ctx, userId, year, month)
	if err != nil {
		return nil, err
	}
	return history, err
}
