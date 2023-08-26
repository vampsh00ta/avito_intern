package transport

import (
	"avito/internal/response"
	"avito/internal/service"
	"avito/internal/transport/model"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"go.uber.org/zap"
	"net/http"
	"strconv"
)

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
func (h HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateOrDeleteUser
	validate := validator.New()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("CreateUser", "mehthod", r.Method, "error", err)
		//val := err.(json.InvalidUnmarshalError)
		response.ReturnError(w, r, err)
		return

	}
	err := validate.Struct(req)
	if err != nil {
		errArr := err.(validator.ValidationErrors)
		h.log.Error("transport:CreateUser ", zap.String("data", err.Error()))
		response.ReturnError(w, r, errArr[0])
		return
	}
	if err := h.service.CreateUser(r.Context(), req.Username); err != nil {
		h.log.Errorw("CreateUser", "method", r.Method, "error", err)
		response.ReturnError(w, r, err)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("CreateUser", "method", r.Method, "status", 201)
	response.ReturnOk(w, r)
	return

}

func (h HttpServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req model.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(req)
		http.Error(w, "validation error", http.StatusBadRequest)
		h.log.Error("transport:AddSegment ", zap.String("data", err.Error()))

		response.ReturnError(w, r, err)
		return

	}
	if err := h.service.DeleteUser(r.Context(), req.Id); err != nil {
		response.ReturnError(w, r, err)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("DeleteUser", "method", r.Method, "status", 201)
	response.ReturnOk(w, r)
	return
}

func (h HttpServer) AddSegment(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateOrDeleteSegment

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(req)
		http.Error(w, "validation error", http.StatusBadRequest)
		h.log.Error("transport:AddSegment ", zap.String("data", err.Error()))

		response.ReturnError(w, r, err)
		return

	}
	if err := h.service.CreateSegment(r.Context(), req.Slug); err != nil {
		response.ReturnError(w, r, err)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("AddSegment", "method", r.Method, "status", 201)

	response.ReturnOk(w, r)
	return
}

func (h HttpServer) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateOrDeleteSegment

	h.log.Info("DeleteSegment")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "validation error", http.StatusBadRequest)
		h.log.Error("transport:DeleteSegment ", zap.String("data", err.Error()))

		response.ReturnError(w, r, err)
		return

	}

	if err := h.service.DeleteSegment(r.Context(), req.Slug); err != nil {
		response.ReturnError(w, r, err)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteSegment", "method", r.Method, "status", 200)
	response.ReturnOk(w, r)
	return
}

func (h HttpServer) GetUsersSegments(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GetUsersSegments ")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ReturnError(w, r, err)
		return
	}

	segments, err := h.service.GetUsersSegments(r.Context(), id)

	if err != nil {
		response.ReturnError(w, r, err)
		return
	}
	h.log.Infow("GetUsersSegments", "method", r.Method, "status", 200)

	response.ReturnOkData(w, r, model.ResponseGetUsersSegments{model.User{id}, segments})
	return

}
func (h HttpServer) AddSegmentsToUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestAddOrDeleteSegmentsToUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "validation error", http.StatusBadRequest)
		h.log.Error("transport:AddSegmentsToUser ", zap.String("data", err.Error()))

		response.ReturnError(w, r, err)
		return

	}

	if err := h.service.AddSegmentsToUser(r.Context(), req); err != nil {
		response.ReturnError(w, r, err)
		return
	}

	h.log.Infow("AddSegmentsToUser", "method", r.Method, "status", 200)
	response.ReturnOk(w, r)
	return
}

func (h HttpServer) DeleteSegmentsFromUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestAddOrDeleteSegmentsToUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "validation error", http.StatusBadRequest)
		h.log.Error("transport:DeleteSegmentsFromUser ", zap.String("data", err.Error()))
		response.ReturnError(w, r, err)
		return

	}

	if err := h.service.DeleteSegmentsFromUser(r.Context(), req); err != nil {
		response.ReturnError(w, r, err)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("AddSegmentsToUser", "method", r.Method, "status", 200)

	response.ReturnOk(w, r)
	return
}

func (h HttpServer) GetHistory(w http.ResponseWriter, r *http.Request) {
	var req model.RequestGetHistory
	var decoder = schema.NewDecoder()
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		h.log.Error("transport:GetHistory ", zap.String("data", err.Error()))
		response.ReturnError(w, r, err)
		return
	}
	history, err := h.service.GetHistory(r.Context(), req.UserID, req.Year, req.Month)
	if err != nil {
		h.log.Error("transport:GetHistory ", zap.String("data", err.Error()))
		response.ReturnError(w, r, err)
		return
	}
	response.ReturnOkData(w, r, history)
	return
}
