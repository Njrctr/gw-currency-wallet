package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	exchanger_grpc "github.com/Njrctr/gw-currency-wallet/internal/clients/exchanger"
	"github.com/Njrctr/gw-currency-wallet/internal/config"
	"github.com/Njrctr/gw-currency-wallet/internal/handlers"
	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository"
	"github.com/Njrctr/gw-currency-wallet/internal/repository/postgres"
	"github.com/Njrctr/gw-currency-wallet/internal/service"
	"github.com/sirupsen/logrus"
)

func Run() {

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("Ошибка инициализации конфига: %s", err.Error())
	}
	
	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		logrus.Fatalf("Ошибка инициализации БД: %s", err.Error())
	}

	exchangeClient, err := exchanger_grpc.NewGRPCClient(
		context.Background(),
		cfg.Client.Address,
	)
	if err != nil {
		logrus.Fatalf("Ошибка инициализации grpc Клиента: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services, exchangeClient, cfg.App.TokenTTL, cfg.App.CacheTTL)
	server := new(models.Server)
	logrus.Print("Try to start server on port: ", cfg.App.Port)

	go func() {
		if err := server.Run(cfg.App.Port, handlers.InitRouters()); err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка запуска сервера: %s", err.Error())
		}
	}()
	logrus.Print("Walleter Backend Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Printf("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Ошибка остановки сервера: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("Ошибка закрытия БД: %s", err.Error())
	}
}
