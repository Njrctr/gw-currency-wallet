package handlers

import (
	"bytes"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/service"
	mock_service "github.com/Njrctr/gw-currency-wallet/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_GetBalance(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, userId int)

	testTable := []struct {
		name                string
		userId              int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "success",
			userId: 1,
			mockBehavior: func(s *mock_service.MockWallet, userId int) {
				s.EXPECT().GetWallet(gomock.Any(), userId).Return(models.Wallet{Balance: models.Balance{
					Usd: 0,
					Rub: 0,
					Eur: 0,
				}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"balance":{"USD":0,"RUB":0,"EUR":0}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mock_service.NewMockWallet(c)
			testCase.mockBehavior(wallet, testCase.userId)

			services := &service.Service{Wallet: wallet}
			handler := NewHandler(services, nil, 0, 0, nil)

			// Test Server
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			r.GET("/balance", func(c *gin.Context) {
				c.Set(userCtx, testCase.userId)
			}, handler.GetBalance)

			// Test HTTP Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/balance", nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_Deposit(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, userId int, input models.EditWallet)

	testTable := []struct {
		name                string
		userId              int
		inputBody           string
		inputEditWallet     models.EditWallet
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "succes",
			userId:    1,
			inputBody: `{"amount":10,"currency":"USD"}`,
			inputEditWallet: models.EditWallet{
				Amount:        float64(10),
				Currency:      "USD",
				OperationType: "DEPOSIT",
			},
			mockBehavior: func(s *mock_service.MockWallet, userId int, input models.EditWallet) {
				s.EXPECT().WithdrawOrDeposit(gomock.Any(), userId, input).Return(models.Balance{
					Usd: 10,
					Rub: 0,
					Eur: 0,
				},
					nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"Account topped up successfully","new_balance":{"USD":10,"RUB":0,"EUR":0}}`,
		},
		{
			name:                "invalid input 1",
			userId:              1,
			inputBody:           `{"currency":"USD"}`,
			mockBehavior:        func(s *mock_service.MockWallet, userId int, input models.EditWallet) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "invalid amount",
			userId:    1,
			inputBody: `{"amount":-10,"currency":"USD"}`,

			mockBehavior:        func(s *mock_service.MockWallet, userId int, input models.EditWallet) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid amount"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mock_service.NewMockWallet(c)
			testCase.mockBehavior(wallet, testCase.userId, testCase.inputEditWallet)

			services := &service.Service{Wallet: wallet}
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			handler := NewHandler(services, nil, 0, 0, log)

			// Test Server
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			r.POST("/deposit", func(c *gin.Context) {
				c.Set(userCtx, testCase.userId)
			}, handler.Deposit)

			// Test HTTP Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/deposit",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})

	}
}

func TestHandler_Withdraw(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, userId int, input models.EditWallet)

	testTable := []struct {
		name                string
		userId              int
		inputBody           string
		inputEditWallet     models.EditWallet
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "succes",
			userId:    1,
			inputBody: `{"amount":10,"currency":"USD"}`,
			inputEditWallet: models.EditWallet{
				Amount:        float64(10),
				Currency:      "USD",
				OperationType: "WITHDRAW",
			},
			mockBehavior: func(s *mock_service.MockWallet, userId int, input models.EditWallet) {
				s.EXPECT().WithdrawOrDeposit(gomock.Any(), userId, input).Return(models.Balance{
					Usd: 10,
					Rub: 0,
					Eur: 0,
				},
					nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"Withdrawal successful","new_balance":{"USD":10,"RUB":0,"EUR":0}}`,
		},
		{
			name:                "invalid input 1",
			userId:              1,
			inputBody:           `{"currency":"USD"}`,
			mockBehavior:        func(s *mock_service.MockWallet, userId int, input models.EditWallet) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "invalid amount",
			userId:    1,
			inputBody: `{"amount":-10,"currency":"USD"}`,

			mockBehavior:        func(s *mock_service.MockWallet, userId int, input models.EditWallet) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid amount"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mock_service.NewMockWallet(c)
			testCase.mockBehavior(wallet, testCase.userId, testCase.inputEditWallet)

			services := &service.Service{Wallet: wallet}
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			handler := NewHandler(services, nil, 0, 0, log)

			// Test Server
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			r.POST("/withdraw", func(c *gin.Context) {
				c.Set(userCtx, testCase.userId)
			}, handler.Withdraw)

			// Test HTTP Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/withdraw",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})

	}
}
