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

func (s service) AddSegmentsToUser(ctx context.Context, userId int, segments ...*Segment_AddSegmentsToUser) error {
	var dbSegments []*rep.Segment
	for _, segment := range segments {
		if segment.Expire != nil {
			days := time.Duration(segment.Expire.Days) * time.Hour * 24
			hours := time.Duration(segment.Expire.Hours) * time.Hour
			minutes := time.Duration(segment.Expire.Minutes) * time.Minute
			timeEnds := time.Now().Add(days + hours + minutes)
			s.ttl.Set(userId, segment.Slug, timeEnds)
		}
		dbSegments = append(dbSegments, &segment.Segment)

	}
	if err := s.rep.AddSegmentsToUser(ctx, userId, dbSegments...); err != nil {
		return err
	}

	return nil
}

func (s service) GetUsersSegments(ctx context.Context, userId int) (*[]rep.Segment, error) {
	return s.rep.GetUsersSegments(ctx, userId)
}
func (s service) DeleteSegmentsFromUser(ctx context.Context, userId int, segments ...*Segment_DeleteSegmentsFromUser) (err error) {
	var keysToDelete []string
	var dbSegments []*rep.Segment

	for _, segment := range segments {
		key := strconv.Itoa(userId) + ":" + segment.Slug
		keysToDelete = append(keysToDelete, key)
		dbSegments = append(dbSegments, &segment.Segment)
	}
	if err := s.rep.DeleteSegmentsFromUser(ctx, userId, dbSegments...); err != nil {
		return err
	}
	s.ttl.Delete(keysToDelete...)
	return nil
}
func (s service) CreateSegmentPercent(ctx context.Context, segment Segment_CreateSegment) (*[]User_CreateSegment, error) {
	slugId, err := s.rep.CreateSegment(ctx, segment.Segment)
	if err != nil {
		return nil, err
	}
	segment.Id = slugId
	if segment.UserPercent == 0 {
		return nil, err
	}
	userIds, err := s.rep.GetUserIds(ctx)
	if err != nil {
		return nil, err
	}
	shuffledUserIds, err := shuffleUsers(userIds, segment.UserPercent)

	if err != nil {

	}
	if err := s.rep.AddSlugIdToUsers(ctx, segment.Segment, shuffledUserIds...); err != nil {
		return nil, err
	}
	var users []User_CreateSegment

	for _, id := range shuffledUserIds {
		users = append(users, User_CreateSegment{rep.User{id}})
	}

	return &users, nil
}
func (s service) CreateSegment(ctx context.Context, segment Segment_CreateSegment) error {

	_, err := s.rep.CreateSegment(ctx, segment.Segment)
	if err != nil {
		return err
	}
	return nil

}
func (s service) DeleteSegment(ctx context.Context, segment Segment_DeleteSegment) error {
	return s.rep.DeleteSegment(ctx, segment.Segment)

}

func (s service) GetHistory(ctx context.Context, userId, year, month int) (*[]rep.HistoryRow, error) {
	var history *[]rep.HistoryRow
	var err error
	if userId != 0 {
		history, err = s.rep.GetHistoryById(ctx, userId, year, month)
	} else {
		history, err = s.rep.GetHistoryAll(ctx, year, month)
	}
	if err != nil {
		return nil, err
	}
	return history, err
}
