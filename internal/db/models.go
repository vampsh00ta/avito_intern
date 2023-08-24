package repository

type Segment struct {
	Id   int    `json:"id" db:"id"`
	Slug string `json:"slug" db:"slug"`
}
