package model

import db "avito/internal/db"

type ResponseGetUsersSegments struct {
	User
	Segments []db.Segment `json:"segments"`
}

type ResponseGetHistory struct {
	db.HistoryRow
}
