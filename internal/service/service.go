package service

import (
	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository"
)

type Authorization interface {
	CreateUser(login models.User) error
	GenerateJWTToken(user models.UserLogin, tokenTTL int) (string, error)
	ParseToken(token string) (int, error)
}

type Wallet interface {
	GetWallet(userId int) (models.Wallet, error)
	WithdrawOrDeposit(userId int, operation models.EditWallet) (models.Balance, error)
}

type Service struct {
	Authorization
	Wallet
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
		Wallet:        newWalletService(repos.Wallet),
	}
}
