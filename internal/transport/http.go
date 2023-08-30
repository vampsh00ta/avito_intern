package transport

import (
	"avito/internal/service"
	"avito/internal/transport/dto"
	"avito/internal/transport/httpresponse"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

var validate = validator.New()
var decoder = schema.NewDecoder()

// PingExample godoc
//
//	@Summary		CreateUser
//	@Description	Создает пользователя
//	@Tags			User
//	@Accept			json
//	@Param  		username   body      dto.RequestCreateUser  true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/user/new [post]
func (h HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		h.log.Infow("CreateUser", "method", r.Method, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	err := validate.Struct(req)
	if err != nil {
		h.log.Infow("CreateUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(r.Context(), req.Username); err != nil {
		h.log.Infow("CreateUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("CreateUser", "method", r.Method, "status", 201)
	httpresponse.ReturnOk(w)
	return

}

// PingExample godoc
//
//	@Summary		DeleteUser
//	@Description	Создает пользователя
//	@Tags			User
//	@Accept			json
//	@Param  		username   body      dto.RequestDeleteUser  true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/user [delete]
func (h HttpServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("DeleteUser ", "method", r.Method, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	err := validate.Struct(req)
	if err != nil {
		h.log.Infow("DeleteUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteUser(r.Context(), req.Id); err != nil {
		h.log.Infow("DeleteUser", "method", r.Method, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("DeleteUser", "method", r.Method, "status", 201)
	httpresponse.ReturnOk(w)
	return
}

// PingExample godoc
//
//	@Summary		CreateSegment
//	@Description	Создает сегмент
//	@Tags			Segment
//	@Accept			json
//
// @Param  		slug   body      dto.RequestCreateSegment  true  "Слэт"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{string}	string	"ok"
// @Failure		404	{string}	string	"ok"
// @Failure		500	{string}	string	"ok"
// @Router			/segment/new [post]
func (h HttpServer) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCreateSegment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("CreateSegment ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}

	if err := validate.Struct(req); err != nil {
		h.log.Infow("CreateSegment", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}

	var res *[]service.User_CreateSegment
	var err error
	if req.UserPercent == 0 {
		err = h.service.CreateSegment(r.Context(), req.Segment_CreateSegment)
	} else {
		res, err = h.service.CreateSegmentPercent(r.Context(), req.Segment_CreateSegment)

	}
	if err != nil {
		h.log.Infow("CreateSegment", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	h.log.Infow("CreateSegment", "method", r.Method, "status", 201)
	if req.UserPercent == 0 {
		httpresponse.ReturnOk(w)

	} else {
		httpresponse.ReturnOkData(w, res)

	}
	return

}

// PingExample godoc
//
//	@Summary		DeleteSegment
//	@Description	Удаляет сегмент
//	@Tags			Segment
//	@Accept			json
//	@Param  		slug   body      dto.RequestDeleteSegment  true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/segment [delete]
func (h HttpServer) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteSegment

	h.log.Info("DeleteSegment")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("DeleteSegment ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("DeleteSegment ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSegment(r.Context(), req.Segment_DeleteSegment); err != nil {
		h.log.Infow("DeleteSegment ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteSegment", "method", r.Method, "status", 200)
	httpresponse.ReturnOk(w)
	return
}

// PingExample godoc
//
//	@Summary		GetUsersSegments
//	@Description	Возвращает сегметы пользователя
//	@Tags			User
//	@Accept			json
//	@Param  		id   path  int      true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/user/segments/{id} [get]
func (h HttpServer) GetUsersSegments(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Infow("GetUsersSegments ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, errors.New("validation error"), http.StatusBadRequest)
		return
	}

	segments, err := h.service.GetUsersSegments(r.Context(), id)

	if err != nil {
		h.log.Infow("GetUsersSegments ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	h.log.Infow("GetUsersSegments", "method", r.Method, "status", 200)
	httpresponse.ReturnOkData(w, dto.ResponseGetUsersSegments{dto.User{id}, segments})
	return

}

// PingExample godoc
//
//	@Summary		AddSegmentsToUser
//	@Description	Добавляет сегменты пользователю
//	@Tags			User
//	@Accept			json
//	@Param  		id    segments  body  dto.RequestAddSegmentsToUser    true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/user/segments/new [post]
func (h HttpServer) AddSegmentsToUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestAddSegmentsToUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("AddSegmentsToUser ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("AddSegmentsToUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.AddSegmentsToUser(r.Context(), req.Id, req.Segments...); err != nil {
		h.log.Infow("AddSegmentsToUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}

	h.log.Infow("AddSegmentsToUser", "method", r.Method, "status", 200)
	httpresponse.ReturnOk(w)
	return
}

// PingExample godoc
//
//	@Summary		DeleteSegmentsFromUser
//	@Description	Удаляет сегменты пользователя
//	@Tags			User
//	@Accept			json
//	@Param  		id    segments  body  dto.RequestDeleteSegmentsFromUser    true  "Слэт"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/user/segments [delete]
func (h HttpServer) DeleteSegmentsFromUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteSegmentsFromUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("DeleteSegmentsFromUser ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("DeleteSegmentsFromUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSegmentsFromUser(r.Context(), req.Id, req.Segments...); err != nil {
		h.log.Infow("DeleteSegmentsFromUser", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteSegmentsFromUser", "method", r.Method, "status", 200)

	httpresponse.ReturnOk(w)
	return
}

// PingExample godoc
//
//	@Summary		GetHistory
//	@Description	Возвращает историю добавления/удаления сегментов пользователю
//	@Tags			History
//	@Accept			json
//	@Param  		user_id query int  false "Id пользователя"
//	@Param  		month  query  string  true "Месяц"
//	@Param  		year  query  string true "Год"
//	@Produce		json
//	@Success		200 {object}  httpresponse.Response
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/history [get]
func (h HttpServer) GetHistory(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestGetHistory
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		h.log.Infow("GetHistory ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Infow("GetHistory ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	history, err := h.service.GetHistory(r.Context(), req.UserID, req.Year, req.Month)
	if err != nil {
		h.log.Infow("GetHistory ", "method", r.Method, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	gocsv.Marshal(history, w)
	h.log.Infow("GetHistory ", "method", r.Method, "status", 200)

	return
}
