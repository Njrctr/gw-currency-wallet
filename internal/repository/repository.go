package repository

import (
	"context"
	"log/slog"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) error
	GetUser(user models.UserLogin) (models.User, error)
}

type Wallet interface {
	GetWallet(ctx context.Context, userId int) (models.Wallet, error)
	EditBalance(ctx context.Context, userId int, operation models.EditWallet) (models.Balance, error)
	Transfer(ctx context.Context, input models.TransferOperation) (models.Balance, error)
}

type Repository struct {
	Authorization
	Wallet
	log *slog.Logger
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgresRepo(db, log),
		Wallet:        postgres.NewWalletPostgresRepo(db, log),
		log:           log,
	}
}
