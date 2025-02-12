package handlers

import (
	"bytes"
	"errors"
	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/service"
	mock_service "github.com/Njrctr/gw-currency-wallet/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_Registration(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "success",
			inputBody: `{"email":"test@test.com","password":"test","username":"test"}`,
			inputUser: models.User{
				Email: "test@test.com",
				UserLogin: models.UserLogin{
					Username: "test",
					Password: "test",
				},
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"User registered successfully"}`,
		},
		{
			name:                "empty fields",
			inputBody:           `{"email":"test@test.com","password":"test"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "server error",
			inputBody: `{"email":"test@test.com","password":"test","username":"test"}`,
			inputUser: models.User{
				Email: "test@test.com",
				UserLogin: models.UserLogin{
					Username: "test",
					Password: "test",
				},
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(errors.New("Service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"Service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			handler := NewHandler(services, nil, 0, 0, log)

			// Test Server
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/register", handler.Registration)

			// Test HTTP Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())

		})
	}
}
