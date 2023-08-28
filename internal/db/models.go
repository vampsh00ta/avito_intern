package repository

import "time"

type Segment struct {
	Slug       string  `json:"slug" db:"slug" validate:"required" csv:"slug"`
	Expire     *Expire `json:"expire,omitempty" db:"-"`
	RandomSeed int     `json:"random,omitempty" db:"-"`

	//"2015-07-05T22:16:18Z
}

type Expire struct {
	Days    int `json:"days,omitempty" db:"-"`
	Hours   int `json:"hours,omitempty" db:"-"`
	Minutes int `json:"minutes,omitempty" db:"-"`
}
type User struct {
	Id int `json:"id" db:"id"`
	//Username string `json:"username" db:"username"`
}
type HistoryRow struct {
	UserId int `json:"user_id" db:"user_id" csv:"user_id"`
	Segment
	Operation  string    `json:"operation" id:"operation" csv:"operation"`
	UpdateTime time.Time `json:"update_time" db:"update_time" csv:"update_time"`
}
