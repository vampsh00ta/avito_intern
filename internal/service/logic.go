package service

import (
	rep "avito/internal/db"
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
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
	//var ttlSlugs []string
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
	fmt.Println(dbSegments)
	if err := s.rep.AddSegmentsToUser(ctx, userId, dbSegments...); err != nil {
		return err
	}
	//if err := s.ttl.SetTTL(ctx, ttlSlugs...); err != nil {
	//	return err
	//}
	return nil
}

func (s service) GetUsersSegments(ctx context.Context, userId int) ([]rep.Segment, error) {
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
func (s service) CreateSegmentPercent(ctx context.Context, segment Segment_CreateSegment) (*[]User_CreateSegmentPercent, error) {
	slugId, err := s.rep.CreateSegment(ctx, segment.Segment)
	if err != nil {
		return nil, err
	}
	segment.Id = slugId
	fmt.Println(segment.UserPercent)
	if segment.UserPercent == 0 {
		return nil, err
	}
	userIds, err := s.rep.GetUserIds(ctx)
	fmt.Println(userIds)
	if err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		return nil, errors.New("null users")
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIds), func(i, j int) { userIds[i], userIds[j] = userIds[j], userIds[i] })
	var percent float64 = float64(segment.UserPercent) / 100
	randomCount := math.Ceil(float64(len(userIds)) * float64(percent))
	shuffledUserIds := userIds[0:int(randomCount)]

	if err := s.rep.AddSlugIdToUsers(ctx, segment.Segment, shuffledUserIds...); err != nil {
		return nil, err
	}
	var users []rep.User

	for _, id := range shuffledUserIds {
		users = append(users, rep.User{id})
	}

	data, err := TypeConverter[[]User_CreateSegmentPercent](&users)
	if err != nil {
		return nil, err
	}
	return data, nil
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
	history, err := s.rep.GetHistory(ctx, userId, year, month)
	if err != nil {
		return nil, err
	}
	return history, err
}
