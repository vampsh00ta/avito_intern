package service

import (
	rep "avito/internal/db"
	"context"
	"strconv"
	"time"
)

func (s service) CreateUser(ctx context.Context, username string) error {

	return s.rep.CreateUser(ctx, username)
}

func (s service) DeleteUser(ctx context.Context, userId int) error {

	return s.rep.DeleteUser(ctx, userId)
}

func (s service) AddSegmentsToUser(ctx context.Context, userId int, segments ...rep.Segment) error {
	var slugs []any
	var ttlSlugs []string
	for _, segment := range segments {
		slugs = append(slugs, segment.Slug)
		if segment.Expire != nil {
			days := time.Duration(segment.Expire.Days) * time.Hour * 24
			hours := time.Duration(segment.Expire.Hours) * time.Hour
			minutes := time.Duration(segment.Expire.Minutes) * time.Minute

			timeEnds := time.Now().Add(days + hours + minutes)
			s.ttl.Collect(&ttlSlugs, userId, segment.Slug, timeEnds)
		}

	}

	if err := s.rep.AddSegmentsToUser(ctx, userId, slugs...); err != nil {
		return err
	}
	if err := s.ttl.SetTTL(ctx, ttlSlugs...); err != nil {
		return err
	}
	return nil
}

func (s service) GetUsersSegments(ctx context.Context, userId int) ([]rep.Segment, error) {
	return s.rep.GetUsersSegments(ctx, userId)
}
func (s service) DeleteSegmentsFromUser(ctx context.Context, userId int, slugs ...any) (err error) {
	var slugsToDelete []string
	for _, slug := range slugs {
		redisKey := strconv.Itoa(userId) + ":" + slug.(string)
		slugsToDelete = append(slugsToDelete, redisKey)

	}
	if err := s.rep.DeleteSegmentsFromUser(ctx, userId, slugs...); err != nil {
		return err
	}
	if err := s.ttl.DelUsersSegments(ctx, slugsToDelete...); err != nil {
		return err
	}
	return nil
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
