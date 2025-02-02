package handlers

import (
	"net/http"
	"reflect"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetBalance
// @Summary Get Wallet Balance
// @Security ApiKeyAuth
// @Tags Wallets
// @Description get wallet balance
// @ID get-wallet-balance
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Wallet
// @Failure 500 {object} errorResponse
// @Router /api/v1/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		return
	}
	wallet, err := h.services.GetWallet(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, wallet)
}

// Deposit
// @Summary Wallet Deposit
// @Security ApiKeyAuth
// @Tags Wallets
// @Description wallet deposit
// @ID wallet-deposit
// @Accept  json
// @Produce  json
// @Param input body models.EditWallet true "deposit input"
// @Success 200 {object} newBalanceResponse
// @Failure 400 {object} errorResponse
// @Router /api/v1/wallet/deposit [post]
func (h *Handler) Deposit(c *gin.Context) {
	var input models.EditWallet

	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	input.OperationType = "DEPOSIT"
	logrus.Println(input)
	logrus.Println(reflect.TypeOf(input.Amount))

	newBalance, err := h.services.WithdrawOrDeposit(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newBalanceResponse{
		Message:    "Account topped up successfully",
		NewBalance: newBalance,
	})
}

type newBalanceResponse struct {
	Message    string         `json:"message"`
	NewBalance models.Balance `json:"new_balance"`
}

// Withdraw
// @Summary Wallet Withdraw
// @Security ApiKeyAuth
// @Tags Wallets
// @Description wallet withdraw
// @ID wallet-withdraw
// @Accept  json
// @Produce  json
// @Param input body models.EditWallet true "withdraw input"
// @Success 200 {object} newBalanceResponse
// @Failure 400 {object} errorResponse
// @Router /api/v1/wallet/withdraw [post]
func (h *Handler) Withdraw(c *gin.Context) {
	var input models.EditWallet

	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	input.OperationType = "WITHDRAW"
	logrus.Println(input)
	logrus.Println(reflect.TypeOf(input.Amount))

	newBalance, err := h.services.WithdrawOrDeposit(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newBalanceResponse{
		Message:    "Withdrawal successful",
		NewBalance: newBalance,
	})
}

func (h *Handler) GetRates(c *gin.Context) {}
