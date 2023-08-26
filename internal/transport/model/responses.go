package model

type ResponseGetUsersSegments struct {
	User
	Segments []Segment `json:"segments"`
}
