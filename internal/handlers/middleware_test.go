package handlers

import (
	"errors"
	"fmt"
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

func TestHandler_userIdentify(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "success",
			headerName:  "Authorization",
			headerValue: "Bearer " + "token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "1",
		},
		{
			name:                "Empty header",
			headerName:          "Authorization",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"Empty auth header"}`,
		},
		{
			name:                "Invalid header",
			headerName:          "Authorization",
			headerValue:         "Bearer",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"Invalid auth header"}`,
		},
		{
			name:                "Empty token",
			headerName:          "Authorization",
			headerValue:         "Bearer " + "",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"Token is empty"}`,
		},
		{
			name:                "Not Bearer",
			headerName:          "Authorization",
			headerValue:         "NotBearer " + "token",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"Invalid auth header. Should be Bearer!"}`,
		},
		{
			name:        "fail",
			headerName:  "Authorization",
			headerValue: "Bearer " + "token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("failed to parse token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			handler := NewHandler(services, nil, 0, 0, log)

			// Test Server
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.GET("/protected", handler.userIdentify, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Make Request
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
