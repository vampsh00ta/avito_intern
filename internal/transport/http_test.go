package transport

import (
	repository "avito/internal/db"
	"avito/internal/service"
	mock_service "avito/internal/service/mocks"
	"avito/internal/transport/dto"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func LoadLoggerDev() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	return sugar
}
func TestTransport_CreateSegment(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService dto.RequestCreateSegment
		expectedCode int
		f            func(s *mock_service.MockService, segment service.Segment_CreateSegment)
		expectedBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug":"test"}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment: repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 201,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegment(gomock.Any(), segment).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"ok"}`,
		},
		{
			name:      "validation",
			inputBody: `{"slufg":"test"}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment: repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 400,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegment(gomock.Any(), segment).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
		{
			name:      "already exists/server",
			inputBody: `{"slug":"test"}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment: repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 500,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegment(gomock.Any(), segment).Return(errors.New("already exists")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"already exists"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)

			test.f(srvc, test.inputService.Segment_CreateSegment)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("POST").Path("/segment/new").HandlerFunc(transport.CreateSegment)
			req := httptest.NewRequest("POST", "/segment/new", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
func TestTransport_DeleteSegment(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService dto.RequestDeleteSegment
		expectedCode int
		f            func(s *mock_service.MockService, segment service.Segment_DeleteSegment)
		expectedBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug":"test"}`,
			inputService: dto.RequestDeleteSegment{
				service.Segment_DeleteSegment{
					repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 200,
			f: func(s *mock_service.MockService, segment service.Segment_DeleteSegment) {
				s.EXPECT().DeleteSegment(gomock.Any(), segment).Return(nil).AnyTimes()
			},

			expectedBody: `{"status":"ok"}`,
		},
		{
			name:      "validation",
			inputBody: `{"slufg":"test"}`,
			inputService: dto.RequestDeleteSegment{
				service.Segment_DeleteSegment{
					repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 400,
			f: func(s *mock_service.MockService, segment service.Segment_DeleteSegment) {
				s.EXPECT().DeleteSegment(gomock.Any(), segment).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		}, {
			name:      "already exists/server",
			inputBody: `{"slug":"test"}`,
			inputService: dto.RequestDeleteSegment{
				service.Segment_DeleteSegment{
					repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 500,
			f: func(s *mock_service.MockService, segment service.Segment_DeleteSegment) {
				s.EXPECT().DeleteSegment(gomock.Any(), segment).Return(errors.New("smth")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)

			test.f(srvc, test.inputService.Segment_DeleteSegment)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("DELETE").Path("/segment").HandlerFunc(transport.DeleteSegment)
			req := httptest.NewRequest("DELETE", "/segment", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
func TestTransport_CreateSegmentPercent(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService dto.RequestCreateSegment
		expectedCode int
		f            func(s *mock_service.MockService, segment service.Segment_CreateSegment)
		expectedBody string
	}{
		{
			name:      "OK",
			inputBody: `{"slug":"test","user_percent":50}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment:     repository.Segment{Slug: "test"},
					UserPercent: 50,
				},
			},
			expectedCode: 201,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegmentPercent(gomock.Any(), segment).Return(&[]service.User_CreateSegment{
					service.User_CreateSegment{
						repository.User{
							Id: 1,
						},
					},
					service.User_CreateSegment{
						repository.User{
							Id: 2,
						},
					},
				}, nil).AnyTimes()
			},
			expectedBody: `{"status":"ok","response":[{"id":1},{"id":2}]}`,
		},
		{
			name:      "validation",
			inputBody: `{"slufg":"test","user_percent":50}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment: repository.Segment{Slug: "test"},
				},
			},
			expectedCode: 400,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegmentPercent(gomock.Any(), segment).Return(&[]service.User_CreateSegment{
					service.User_CreateSegment{
						repository.User{
							Id: 1,
						},
					},
					service.User_CreateSegment{
						repository.User{
							Id: 2,
						},
					},
				}, nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		}, {
			name:      "already exists/server",
			inputBody: `{"slug":"test","user_percent":50}`,
			inputService: dto.RequestCreateSegment{
				service.Segment_CreateSegment{
					Segment:     repository.Segment{Slug: "test"},
					UserPercent: 50,
				},
			},
			expectedCode: 500,
			f: func(s *mock_service.MockService, segment service.Segment_CreateSegment) {
				s.EXPECT().CreateSegmentPercent(gomock.Any(), segment).Return(nil, errors.New("already exists"))
			},
			expectedBody: `{"status":"error","error":"already exists"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)

			test.f(srvc, test.inputService.Segment_CreateSegment)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("POST").Path("/segment/new").HandlerFunc(transport.CreateSegment)
			req := httptest.NewRequest("POST", "/segment/new", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestTransport_CreateUser(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService string
		expectedCode int
		f            func(s *mock_service.MockService, data string)
		expectedBody string
	}{
		{
			name:         "OK",
			inputBody:    `{"username":"test"}`,
			inputService: "test",
			expectedCode: 201,
			f: func(s *mock_service.MockService, data string) {
				s.EXPECT().CreateUser(gomock.Any(), data).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"ok"}`,
		},
		{
			name:         "already exists/server",
			inputBody:    `{"username":"test"}`,
			inputService: "test",
			expectedCode: 500,
			f: func(s *mock_service.MockService, data string) {
				s.EXPECT().CreateUser(gomock.Any(), data).Return(errors.New("already exists")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"already exists"}`,
		},
		{
			name:         "ERROR",
			inputBody:    `{"dsf":"f"}`,
			inputService: "test",
			expectedCode: 400,
			f: func(s *mock_service.MockService, data string) {
				s.EXPECT().CreateUser(gomock.Any(), data).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)
			test.f(srvc, test.inputService)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("POST").Path("/user/new").HandlerFunc(transport.CreateUser)
			req := httptest.NewRequest("POST", "/user/new", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
func TestTransport_DeleteUser(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService int
		expectedCode int
		f            func(s *mock_service.MockService, data int)
		expectedBody string
	}{
		{
			name:         "OK",
			inputBody:    `{"id":1}`,
			inputService: 1,
			expectedCode: 200,
			f: func(s *mock_service.MockService, data int) {
				s.EXPECT().DeleteUser(gomock.Any(), data).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"ok"}`,
		},
		{
			name:         "validation",
			inputBody:    `{"f":3}`,
			inputService: 1,
			expectedCode: 400,
			f: func(s *mock_service.MockService, data int) {
				s.EXPECT().DeleteUser(gomock.Any(), data).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
		{
			name:         "already exists/server",
			inputBody:    `{"id":3}`,
			inputService: 3,
			expectedCode: 500,
			f: func(s *mock_service.MockService, data int) {
				s.EXPECT().DeleteUser(gomock.Any(), data).Return(errors.New("smth")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)
			test.f(srvc, test.inputService)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("DELETE").Path("/user").HandlerFunc(transport.DeleteUser)
			req := httptest.NewRequest("DELETE", "/user", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestTransport_AddSegmentsToUser(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService dto.RequestAddSegmentsToUser
		expectedCode int
		f            func(s *mock_service.MockService, id int, segments ...any)
		expectedBody string
	}{
		{
			name:      "OK",
			inputBody: `{ "id": 1, "segments": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestAddSegmentsToUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_AddSegmentsToUser{
					&service.Segment_AddSegmentsToUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 200,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().AddSegmentsToUser(gomock.Any(), id, segments...).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"ok"}`,
		},
		{
			name:      "validation",
			inputBody: `{ "id": 1, "segment": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestAddSegmentsToUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_AddSegmentsToUser{
					&service.Segment_AddSegmentsToUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 400,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().AddSegmentsToUser(gomock.Any(), id, segments...).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
		{
			name:      "already exists/server",
			inputBody: `{ "id": 1, "segments": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestAddSegmentsToUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_AddSegmentsToUser{
					&service.Segment_AddSegmentsToUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 500,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().AddSegmentsToUser(gomock.Any(), id, segments...).Return(errors.New("already exists")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"already exists"}`,
		},
		//{
		//	name:      "validation",
		//	inputBody: `{ id": 1, "segments": [ {"slug": "test1" } ] }`,
		//	inputService: dto.RequestAddSegmentsToUser{
		//		User: dto.User{
		//			Id: 1,
		//		},
		//		Segments: []*service.Segment_AddSegmentsToUser{
		//			&service.Segment_AddSegmentsToUser{
		//				Segment: repository.Segment{Slug: "test1"},
		//			},
		//		},
		//	},
		//	expectedCode: 400,
		//	f: func(s *mock_service.MockService, id int, segments ...any) {
		//		s.EXPECT().AddSegmentsToUser(gomock.Any(), id, segments...).Return(errors.New("already exists")).AnyTimes()
		//	},
		//	expectedBody: `{"status":"error","error":"already exists"}`,
		//},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)
			a := []any{}
			for _, b := range test.inputService.Segments {
				a = append(a, b)
			}
			test.f(srvc, test.inputService.Id, a...)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("POST").Path("/user/segments/add").HandlerFunc(transport.AddSegmentsToUser)
			req := httptest.NewRequest("POST", "/user/segments/add", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
func TestTransport_DeleteFromUser(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService dto.RequestDeleteSegmentsFromUser
		expectedCode int
		f            func(s *mock_service.MockService, id int, segments ...any)
		expectedBody string
	}{
		{
			name:      "OK",
			inputBody: `{ "id": 1, "segments": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestDeleteSegmentsFromUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_DeleteSegmentsFromUser{
					&service.Segment_DeleteSegmentsFromUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 200,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().DeleteSegmentsFromUser(gomock.Any(), id, segments...).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"ok"}`,
		},
		{
			name:      "validation",
			inputBody: `{ "id": 1, "sement": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestDeleteSegmentsFromUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_DeleteSegmentsFromUser{
					&service.Segment_DeleteSegmentsFromUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 400,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().DeleteSegmentsFromUser(gomock.Any(), id, segments...).Return(nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
		{
			name:      "already exists/error",
			inputBody: `{ "id": 1, "segments": [ {"slug": "test1" } ] }`,
			inputService: dto.RequestDeleteSegmentsFromUser{
				User: dto.User{
					Id: 1,
				},
				Segments: []*service.Segment_DeleteSegmentsFromUser{
					&service.Segment_DeleteSegmentsFromUser{
						Segment: repository.Segment{Slug: "test1"},
					},
				},
			},
			expectedCode: 500,
			f: func(s *mock_service.MockService, id int, segments ...any) {
				s.EXPECT().DeleteSegmentsFromUser(gomock.Any(), id, segments...).Return(errors.New("smth")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"server error"}`,
		},

		//{
		//	name:         "ERROR",
		//	inputBody:    `{"f":3}`,
		//	inputService: 1,
		//	expectedCode: 400,
		//	f: func(s *mock_service.MockService, data int) {
		//		s.EXPECT().DeleteUser(gomock.Any(), data).Return(nil).AnyTimes()
		//	},
		//	expectedBody: `{"status":"error","error":"validation error"}`,
		//},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)
			a := []any{}
			for _, b := range test.inputService.Segments {
				a = append(a, b)
			}
			test.f(srvc, test.inputService.Id, a...)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Methods("DELETE").Path("/user/segments").HandlerFunc(transport.DeleteSegmentsFromUser)
			req := httptest.NewRequest("DELETE", "/user/segments", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
func TestTransport_GetUsersSegments(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		inputService int
		expectedCode int
		f            func(s *mock_service.MockService, id int)
		expectedBody string
	}{
		{
			name:         "OK",
			inputBody:    "1",
			inputService: 1,
			expectedCode: 200,
			f: func(s *mock_service.MockService, id int) {
				s.EXPECT().GetUsersSegments(gomock.Any(), id).Return(&[]repository.Segment{
					{Slug: "test1"},
					{Slug: "test2"},
				}, nil).AnyTimes()
			},
			expectedBody: `{"status":"ok","response":{"id":1,"segments":[{"slug":"test1"},{"slug":"test2"}]}}`,
		},
		{
			name:         "validation",
			inputBody:    "asdf",
			inputService: 1,
			expectedCode: 400,
			f: func(s *mock_service.MockService, id int) {
				s.EXPECT().GetUsersSegments(gomock.Any(), id).Return(&[]repository.Segment{
					{Slug: "test1"},
					{Slug: "test2"},
				}, nil).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"validation error"}`,
		},
		{
			name:         "already exists/server",
			inputBody:    "1",
			inputService: 1,
			expectedCode: 500,
			f: func(s *mock_service.MockService, id int) {
				s.EXPECT().GetUsersSegments(gomock.Any(), id).Return(nil, errors.New("smth")).AnyTimes()
			},
			expectedBody: `{"status":"error","error":"server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvc := mock_service.NewMockService(ctrl)

			test.f(srvc, test.inputService)

			transport := NewHttpServer(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			fmt.Println("/user/segments/" + test.inputBody)
			router.Methods("GET").Path("/user/segments/{id}").HandlerFunc(transport.GetUsersSegments)
			req := httptest.NewRequest("GET", "/user/segments/"+test.inputBody, bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

//func TestTransport_GetHistory(t *testing.T) {
//	timeNow := time.Now()
//	expectedBody := fmt.Sprintf("user_id,slug,operation,update_time1,test1,insert,%s,test2,delete,%s", timeNow, timeNow.String())
//	tests := []struct {
//		name         string
//		inputBody    dto.RequestGetHistory
//		inputService dto.RequestGetHistory
//		expectedCode int
//		f            func(s *mock_service.MockService, data dto.RequestGetHistory)
//		expectedBody string
//	}{
//		{
//			name: "OK",
//			inputBody: dto.RequestGetHistory{
//				UserID: 1,
//				Month:  8,
//				Year:   2023,
//			},
//			inputService: dto.RequestGetHistory{
//				UserID: 1,
//				Month:  8,
//				Year:   2023,
//			},
//			expectedCode: 200,
//			f: func(s *mock_service.MockService, data dto.RequestGetHistory) {
//				s.EXPECT().GetHistory(gomock.Any(), data.UserID, data.Year, data.Month).Return(&[]repository.HistoryRow{
//					{UserId: 1, Segment: repository.Segment{Slug: "test1"}, Operation: "insert", UpdateTime: timeNow},
//					{UserId: 2, Segment: repository.Segment{Slug: "test2"}, Operation: "delete", UpdateTime: timeNow},
//				}, nil).AnyTimes()
//			},
//			expectedBody: expectedBody,
//		},
//		{
//			name: "validation",
//			inputBody: dto.RequestGetHistory{
//				UserID: 1,
//				Year:   2023,
//			},
//			inputService: dto.RequestGetHistory{
//				UserID: 1,
//				Month:  8,
//				Year:   2023,
//			},
//			expectedCode: 400,
//			f: func(s *mock_service.MockService, data dto.RequestGetHistory) {
//				s.EXPECT().GetHistory(gomock.Any(), data.UserID, data.Year, data.Month).Return(&[]repository.HistoryRow{
//					{UserId: 1, Segment: repository.Segment{Slug: "test1"}, Operation: "insert", UpdateTime: timeNow},
//					{UserId: 2, Segment: repository.Segment{Slug: "test2"}, Operation: "delete", UpdateTime: timeNow},
//				}, nil).AnyTimes()
//			},
//			expectedBody: `{"status":"error","error":"validation error"}`,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			srvc := mock_service.NewMockService(ctrl)
//
//			test.f(srvc, test.inputService)
//
//			transport := NewHttpServer(srvc, LoadLoggerDev())
//			w := httptest.NewRecorder()
//			router := mux.NewRouter()
//			router.Methods("GET").Path("/history").HandlerFunc(transport.GetHistory)
//			req := httptest.NewRequest("GET",
//				"/history?user_id="+strconv.Itoa(test.inputBody.UserID)+"&month="+strconv.Itoa(test.inputBody.Month)+"&year="+strconv.Itoa(test.inputBody.Year),
//				bytes.NewBufferString(""))
//			router.ServeHTTP(w, req)
//			assert.Equal(t, test.expectedCode, w.Code)
//			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
//		})
//	}
//}
