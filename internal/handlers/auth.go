package handlers

import (
	"fmt"
	"net/http"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/gin-gonic/gin"
)

// Registration
// @Summary Registration
// @Tags Auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.User true "account data"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/register [post]
func (h *Handler) Registration(c *gin.Context) {
	var input models.User

	h.log.With("func", "Handler/Registration")

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Debug(fmt.Sprintf("invalid input body: %v", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.CreateUser(input); err != nil {
		h.log.Debug(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "Service error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

// Login
// @Summary Login
// @Tags Auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body models.UserLogin true "login data"
// @Success 200 {string} string "token"
// @Failure 401 {object} errorResponse
// @Router /api/v1/login [post]
func (h *Handler) Login(c *gin.Context) {
	var input models.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.GenerateJWTToken(input, h.tokenTTL)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
