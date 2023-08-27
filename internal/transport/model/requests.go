package model

import db "avito/internal/db"

type User struct {
	Id int `json:"id" validate:"required"`
}

type RequestCreateOrDeleteUser struct {
	Username string `json:"username" validate:"required"`
}
type RequestCreateOrDeleteSegment struct {
	db.Segment
}

//	type AddOrDeleteSegmentsToUser_Segment struct {
//		Segments []db.Segment `json:"slugs" validate:"required"`
//		Expire int    `json:"expire,omitempty" `
//
// }
type RequestAddOrDeleteSegmentsToUser struct {
	User
	Segments []db.Segment `json:"slugs" validate:"required"`
}

type RequestGetHistory struct {
	UserID int `json:"user_id" validate:"required" schema:"user_id"`
	Month  int `json:"month" validate:"required" schema:"month"`
	Year   int `json:"year" validate:"required" schema:"year"`
}
