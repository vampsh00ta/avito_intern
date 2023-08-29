package dto

import (
	"avito/internal/service"
)

type User struct {
	Id int `json:"id" validate:"required"`
}

type RequestCreateUser struct {
	Username string `json:"username" validate:"required"`
}
type RequestDeleteUser struct {
	User
}
type RequestCreateSegment struct {
	service.Segment_CreateSegment `validate:"required"`
}
type RequestDeleteSegment struct {
	service.Segment_DeleteSegment
}

type RequestAddSegmentsToUser struct {
	User
	Segments []*service.Segment_AddSegmentsToUser `json:"segments" validate:"required"`
}
type RequestDeleteSegmentsFromUser struct {
	User
	Segments []*service.Segment_DeleteSegmentsFromUser `json:"segments" validate:"required"`
}

type RequestGetHistory struct {
	UserID int `json:"user_id" validate:"required" schema:"user_id"`
	Month  int `json:"month" validate:"required" schema:"month"`
	Year   int `json:"year" validate:"required" schema:"year"`
}
