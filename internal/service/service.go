package service

import (
	"context"
	"log/slog"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository"
)

//go:generate mockgen -source=./service.go -destination=./mocks/service.go

type Authorization interface {
	CreateUser(login models.User) error
	GenerateJWTToken(user models.UserLogin, tokenTTL int) (string, error)
	ParseToken(token string) (int, error)
}

type Wallet interface {
	GetWallet(ctx context.Context, userId int) (models.Wallet, error)
	WithdrawOrDeposit(ctx context.Context, userId int, operation models.EditWallet) (models.Balance, error)
	Transfer(ctx context.Context, input models.TransferOperation) (models.Balance, error)
}

type Service struct {
	Authorization
	Wallet
	log *slog.Logger
}

func NewService(repos *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
		Wallet:        newWalletService(repos.Wallet),
		log:           log,
	}
}
