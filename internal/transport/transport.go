package transport

import (
	"avito/internal/service"
	"go.uber.org/zap"
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
type HttpServer struct {
	service service.Service
	log     *zap.SugaredLogger
}

func NewHttpServer(service service.Service, logger *zap.SugaredLogger) Transport {
	return HttpServer{
		service: service,
		log:     logger,
	}
}
