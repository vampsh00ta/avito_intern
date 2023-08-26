package repository

import "time"

type Segment struct {
	Slug string `json:"slug" db:"slug" validate:"required"`
}
type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}
type HistoryRow struct {
	UserId int `json:"user_id" db:"user_id" `
	Segment
	Operation  string    `json:"operation" id:"operation"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}
