package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	wallet, err := h.services.GetWallet(c, userId)
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

	newBalance, err := h.services.WithdrawOrDeposit(c, userId, input)
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

	newBalance, err := h.services.WithdrawOrDeposit(c, userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newBalanceResponse{
		Message:    "Withdrawal successful",
		NewBalance: newBalance,
	})
}

// GetRates
// @Summary Get Rates
// @Security ApiKeyAuth
// @Tags Wallets
// @Description get rates
// @ID get-rates
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Rates
// @Failure 500 {object} errorResponse
// @Router /api/v1/exchange/rates [get]
func (h *Handler) GetRates(c *gin.Context) {

	// Проверяем не устарели ли значения в хэше
	if h.cache.LastUpdate < time.Now().Unix() {
		rates, err := h.exchanges.GetExchangeRates(c)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		defer h.cache.Set(rates.Rates)

		c.JSON(http.StatusOK, gin.H{
			"rates": rates.Rates,
		})
		return
	}
	// Берём значение из хэша
	rates := h.cache.Get()
	c.JSON(http.StatusOK, gin.H{
		"rates": rates,
	})
}

// Exchange
// @Summary Exchange
// @Security ApiKeyAuth
// @Tags Wallets
// @Description Exchange
// @ID exchange
// @Accept  json
// @Produce  json
// @Param input body models.ExchangeRequest true "Exchange input"
// @Success 200 {object} exchangeResponse
// @Failure 400 {object} errorResponse
// @Router /api/v1/exchange [post]
func (h *Handler) Exchange(c *gin.Context) {
	var input models.ExchangeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Insufficient funds or invalid currencies",
		})
		return
	} else if input.From == "" || input.To == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "Insufficient funds or invalid currencies",
		})
		return
	}

	if input.From == input.To {
		c.JSON(http.StatusOK, gin.H{
			"error": "Insufficient funds or invalid currencies",
		})
		return
	}

	rate, err := h.exchanges.GetRateForCurrency(context.Background(), input.From, input.To)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		return
	}

	transferData := models.TransferOperation{
		UserId: userId,
		From:   input.From,
		To:     input.To,
		Amount: input.Amount,
		Rate:   rate.Rate,
	}
	newBalance, err := h.services.Transfer(context.Background(), transferData)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(newBalance)
	c.JSON(http.StatusOK, exchangeResponse{
		Message:    "Exchange successful",
		Amount:     input.Amount,
		NewBalance: newBalance,
	})
}

type exchangeResponse struct {
	Message    string         `json:"message"`
	Amount     float64        `json:"exchanged_amount"`
	NewBalance models.Balance `json:"new_balance"`
}
