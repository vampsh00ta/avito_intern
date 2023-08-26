package model

type Segment struct {
	Slug string `json:"slug" 'db':"slug" validate:"required"`
}
type User struct {
	Id int `json:"id" validate:"required"`
}

type RequestCreateOrDeleteUser struct {
	Username string `json:"username" validate:"required"`
}
type RequestCreateOrDeleteSegment struct {
	Segment
}

type RequestAddOrDeleteSegmentsToUser struct {
	User
	Segments []Segment `json:"slugs" validate:"required"`
}
