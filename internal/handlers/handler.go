package handlers

import (
	_ "github.com/Njrctr/gw-currency-wallet/docs"
	"github.com/Njrctr/gw-currency-wallet/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	tokenTTL int
}

func NewHandler(services *service.Service, tokenTTL int) *Handler {
	return &Handler{
		services: services,
		tokenTTL: tokenTTL,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{
			// AUTH
			v1.POST("/login", h.Login)
			v1.POST("/register", h.Registration)

			// WALLET
			v1.GET("/balance", h.userIdentify, h.GetBalance)

			wallet := v1.Group("/wallet", h.userIdentify)
			{
				wallet.POST("/deposit", h.Deposit)
				wallet.POST("/withdraw", h.Withdraw)
			}

			// EXCHANGE
			exchange := v1.Group("/exchange", h.userIdentify)
			{
				exchange.GET("/rates", h.GetRates)
			}
		}
	}

	return router
}
