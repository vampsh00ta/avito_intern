package transport

import "avito/internal/service"

type HttpServer struct {
	service service.Service
}

func NewHttpServer(service service.Service) HttpServer {
	return HttpServer{
		service: service,
	}
}
