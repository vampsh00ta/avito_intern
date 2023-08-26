package transport

import (
	"net/http"
)

type Transport interface {
	AddSegment(w http.ResponseWriter, r *http.Request)
	DeleteSegment(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	AddSegmentsToUser(w http.ResponseWriter, r *http.Request)
	GetUsersSegments(w http.ResponseWriter, r *http.Request)
	DeleteSegmentsFromUser(w http.ResponseWriter, r *http.Request)
	GetHistory(w http.ResponseWriter, r *http.Request)
}
