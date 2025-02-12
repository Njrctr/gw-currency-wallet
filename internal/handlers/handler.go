package handlers

import (
	"log/slog"

	_ "github.com/Njrctr/gw-currency-wallet/docs"
	exchanger_grpc "github.com/Njrctr/gw-currency-wallet/internal/clients/exchanger"
	"github.com/Njrctr/gw-currency-wallet/internal/service"
	"github.com/Njrctr/gw-currency-wallet/pkg/cache"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services  *service.Service
	exchanges *exchanger_grpc.GRPCClient
	cache     *cache.CacheInMemory
	tokenTTL  int
	log       *slog.Logger
}

func NewHandler(
	services *service.Service,
	grpcapi *exchanger_grpc.GRPCClient,
	tokenTTL int,
	cacheTTL int,
	log *slog.Logger,
) *Handler {
	return &Handler{
		services:  services,
		exchanges: grpcapi,
		cache:     cache.NewCacheInMemory(int64(cacheTTL)),
		tokenTTL:  tokenTTL,
		log:       log,
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
				exchange.POST("", h.Exchange)
				exchange.GET("/rates", h.GetRates)
			}
		}
	}

	return router
}
