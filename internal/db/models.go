package repository

type Segment struct {
	Slug string `json:"slug" db:"slug"`
}
type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}
