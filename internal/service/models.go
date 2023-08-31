package service

import db "avito/internal/db"

type SegmentModel struct {
	Slug string `json:"slug" db:"slug" validate:"required" csv:"slug"`
}
type Segment_CreateSegment struct {
	db.Segment
	UserPercent int `json:"user_percent,omitempty"`
}
type User_CreateSegment struct {
	db.User
}
type Segment_DeleteSegment struct {
	db.Segment
}

type Segment_AddSegmentsToUser struct {
	db.Segment
	Expire *Expire `json:"expire,omitempty" schema:"expire"`
}

type Segment_DeleteSegmentsFromUser struct {
	db.Segment
}

type Segment_GetUserSegment struct {
	db.Segment
}
type Expire struct {
	Days    int `json:"days,omitempty"`
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
}
