package service

import (
	"context"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository"
)

type WalletService struct {
	repo repository.Wallet
}

func newWalletService(repo repository.Wallet) *WalletService {

	return &WalletService{repo: repo}
}

func (s *WalletService) GetWallet(ctx context.Context, userId int) (models.Wallet, error) {

	return s.repo.GetWallet(ctx, userId)
}

func (s *WalletService) WithdrawOrDeposit(ctx context.Context, userId int, operation models.EditWallet) (models.Balance, error) {

	return s.repo.EditBalance(ctx, userId, operation)
}

func (s *WalletService) Transfer(ctx context.Context, input models.TransferOperation) (models.Balance, error) {

	return s.repo.Transfer(ctx, input)
}
