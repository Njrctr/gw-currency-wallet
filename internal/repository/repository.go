package repository

import (
	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/Njrctr/gw-currency-wallet/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) error
	GetUser(user models.UserLogin) (models.User, error)
}

type Wallet interface {
	GetWallet(userId int) (models.Wallet, error)
	EditBalance(userId int, operation models.EditWallet) (models.Balance, error)
}

type Repository struct {
	Authorization
	Wallet
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgresRepo(db),
		Wallet:        postgres.NewWalletPostgresRepo(db),
	}
}
