package service

import (
	mock_service "avito/internal/service/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockService, slug string)
	ctx := context.Background()
	testTable := []struct {
		name           string
		inputBody      string
		mockBehavior   mockBehavior
		expectedResult interface{}
	}{
		{name: "Ok",
			inputBody: "test",
			mockBehavior: func(s *mock_service.MockService, slug string) {
				s.EXPECT().CreateUser(ctx, slug).Return(nil)
			},
			expectedResult: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			user := mock_service.NewMockService(c)
			testCase.mockBehavior(user, testCase.inputBody)
			//services := service{}
		})
	}
}
