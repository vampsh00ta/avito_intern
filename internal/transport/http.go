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

// @Summary		CreateUser
// @Description	Создает пользователя.Принимает имя пользователя.Если такой пользователь уже существует,то выведете соответствующую ошибку.
// @Tags			User
// @Accept			json
// @Param  		username   body      dto.RequestCreateUser  true  "username"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/user/new [post]
func (h HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		h.log.Infow("CreateUser", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	err := validate.Struct(req)
	if err != nil {
		h.log.Infow("CreateUser", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(r.Context(), req.Username); err != nil {
		h.log.Infow("CreateUser", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	h.log.Infow("CreateUser", "method", r.Method, "status", 201)
	httpresponse.ReturnOk(w)
	return

}

// @Summary		DeleteUser
// @Description	Удаляет пользователя.Принимает Id пользователя
// @Tags			User
// @Accept			json
// @Param  		id   body      dto.RequestDeleteUser  true  "id"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/user [delete]
func (h HttpServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("DeleteUser ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	err := validate.Struct(req)
	if err != nil {
		h.log.Infow("DeleteUser", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteUser(r.Context(), req.Id); err != nil {
		h.log.Infow("DeleteUser", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)

		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteUser", "method", r.Method, "status", 201)
	httpresponse.ReturnOk(w)
	return
}

//	@Summary		CreateSegment
//	@Description	Создает сегмент.Принимает slug сегмента и процент пользоватей,которым он присвоится.Если указан user_percent,то добавит созданный сегмент указанному проценту пользователей(округление идет в большую сторону) , добавит запись в history и вернет id пользователей,котором добавили созданный сегмент.Если такой сегмент уже существует,то выведет соответствующую ошибку.
//	@Tags			Segment
//	@Accept			json
//
// @Param  		slug   body      dto.RequestCreateSegment  true  "slug"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/segment/new [post]
func (h HttpServer) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCreateSegment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("CreateSegment ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}

	if err := validate.Struct(req); err != nil {
		h.log.Infow("CreateSegment", "method", r.Method,
			"status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}

	var res *[]service.User_CreateSegment
	var err error
	if req.UserPercent == 0 {
		err = h.service.CreateSegment(r.Context(), req.Segment_CreateSegment)
	} else if req.UserPercent > 0 {
		res, err = h.service.CreateSegmentPercent(r.Context(), req.Segment_CreateSegment)
	} else {
		h.log.Infow("CreateSegment", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}

	if err != nil {
		h.log.Infow("CreateSegment", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
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

// @Summary		DeleteSegment
// @Description	Удаляет сегмент.Принимает slug сегмента
// @Tags			Segment
// @Accept			json
// @Param  		slug   body      dto.RequestDeleteSegment  true  "slug"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/segment [delete]
func (h HttpServer) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteSegment

	h.log.Info("DeleteSegment")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("DeleteSegment ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("DeleteSegment ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSegment(r.Context(), req.Segment_DeleteSegment); err != nil {
		h.log.Infow("DeleteSegment ", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteSegment", "method", r.Method, "status", 200)
	httpresponse.ReturnOk(w)
	return
}

// @Summary		GetUsersSegments
// @Description	Возвращает сегметы пользовател.Принимает id пользователя
// @Tags			User
// @Accept			json
// @Param  		id   path  int      true  "id"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/user/segments/{id} [get]
func (h HttpServer) GetUsersSegments(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Infow("GetUsersSegments ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, errors.New("validation error"), http.StatusBadRequest)
		return
	}

	segments, err := h.service.GetUsersSegments(r.Context(), id)

	if err != nil {
		h.log.Infow("GetUsersSegments ", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	h.log.Infow("GetUsersSegments", "method", r.Method, "status", 200)
	httpresponse.ReturnOkData(w, dto.ResponseGetUsersSegments{dto.User{id}, segments})
	return

}

// @Summary		AddSegmentsToUser
// @Description	Добавляет сегменты пользователю.Принимает id пользователя и сегменты с полями slug и expire.Если в сегментах указаны expire,то добавляет заданным сегментам TTL.Также добавляет запись добавления в history
// @Tags			User
// @Accept			json
// @Param  		id    segments  body  dto.RequestAddSegmentsToUser    true  "id"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/user/segments/add [post]
func (h HttpServer) AddSegmentsToUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestAddSegmentsToUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("AddSegmentsToUser ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("AddSegmentsToUser", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.AddSegmentsToUser(r.Context(), req.Id, req.Segments...); err != nil {
		h.log.Infow("AddSegmentsToUser", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}

	h.log.Infow("AddSegmentsToUser", "method", r.Method, "status", 200)
	httpresponse.ReturnOk(w)
	return
}

// @Summary		DeleteSegmentsFromUser
// @Description	Удаляет сегменты пользователя.Принимает slug сегмента.Если у сегмента был TTL, то удаляет его из кэша.Также добавляет запись удаления в history
// @Tags			User
// @Accept			json
// @Param  		id    segments  body  dto.RequestDeleteSegmentsFromUser    true  "id"
// @Produce		json
// @Success		200 {object}  httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/user/segments [delete]
func (h HttpServer) DeleteSegmentsFromUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestDeleteSegmentsFromUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("DeleteSegmentsFromUser ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return

	}
	if err := validate.Struct(req); err != nil {
		h.log.Infow("DeleteSegmentsFromUser", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSegmentsFromUser(r.Context(), req.Id, req.Segments...); err != nil {
		h.log.Infow("DeleteSegmentsFromUser", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	h.log.Infow("DeleteSegmentsFromUser", "method", r.Method, "status", 200)

	httpresponse.ReturnOk(w)
	return
}

// @Summary		GetHistory
// @Description	Возвращает историю добавления/удаления сегментов  в виде csv файла.Принимает id пользователя и период времени .Если указан user_id , то возвращает историю конкретного пользователя, иначе - все историю заданного периода
// @Tags			History
// @Accept			json
// @Param  		user_id query int  false "user_id"
// @Param  		month  query  string  true "Месяц"
// @Param  		year  query  string true "Год"
// @Produce		json
//
// @Success		200 {object}	httpresponse.Response
// @Failure		400	{object}	httpresponse.Response
// @Failure		404	{object}	httpresponse.Response
// @Failure		500	{object}	httpresponse.Response
// @Router			/history [get]
func (h HttpServer) GetHistory(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestGetHistory
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		h.log.Infow("GetHistory ", "method",
			r.Method, "error", "status", http.StatusBadRequest, err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Infow("GetHistory ", "method",
			r.Method, "status", http.StatusBadRequest, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusBadRequest)
		return
	}
	history, err := h.service.GetHistory(r.Context(), req.UserID, req.Year, req.Month)
	if err != nil {
		h.log.Infow("GetHistory ", "method",
			r.Method, "status", http.StatusInternalServerError, "error", err)
		httpresponse.ReturnError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	gocsv.Marshal(history, w)
	h.log.Infow("GetHistory ", "method", r.Method, "status", 200)

	return
}
