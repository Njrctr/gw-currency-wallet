package postgres

import (
	"fmt"
	"github.com/Njrctr/gw-currency-wallet/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	walletsTable      = "wallets"
	usersTable        = "users"
	usersWalletsTable = "users_wallets"
)

func NewDB(cfg config.ConfigDB) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.DBName, cfg.Username, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
