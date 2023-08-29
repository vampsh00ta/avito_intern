package repository

import "time"

type Segment struct {
	Id   int    `json:"id,omitempty" db:"-"   csv:"-"`
	Slug string `json:"slug" db:"slug" validate:"required" csv:"slug"`
}

type User struct {
	Id int `json:"id,omitempty" db:"id"`
	//Username string `json:"username" db:"username"`
}
type HistoryRow struct {
	UserId int `json:"user_id" db:"user_id" csv:"user_id"`
	Segment
	Operation  string    `json:"operation" id:"operation" csv:"operation"`
	UpdateTime time.Time `json:"update_time" db:"update_time" csv:"update_time"`
}
