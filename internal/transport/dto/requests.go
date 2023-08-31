package dto

import (
	"avito/internal/service"
)

type User struct {
	Id int `json:"id" validate:"required,gt=0,omitempty"`
}

type RequestCreateUser struct {
	Username string `json:"username" validate:"required,gt=0,omitempty"`
}
type RequestDeleteUser struct {
	User `validate:"required,gt=0,omitempty"`
}
type RequestCreateSegment struct {
	service.Segment_CreateSegment `validate:"required,gt=0,omitempty" `
}
type RequestDeleteSegment struct {
	service.Segment_DeleteSegment
}

type RequestAddSegmentsToUser struct {
	User
	Segments []*service.Segment_AddSegmentsToUser `json:"segments" validate:"required,gt=0,omitempty"`
}
type RequestDeleteSegmentsFromUser struct {
	User
	Segments []*service.Segment_DeleteSegmentsFromUser `json:"segments" validate:"required,gt=0,omitempty"`
}

type RequestGetHistory struct {
	UserID int `json:"user_id"  schema:"user_id"`
	Month  int `json:"month" validate:"required"`
	Year   int `json:"year" validate:"required" `
}
