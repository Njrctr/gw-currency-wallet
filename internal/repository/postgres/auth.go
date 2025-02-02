package postgres

import (
	"fmt"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgresRepo struct {
	db *sqlx.DB
}

func NewAuthPostgresRepo(db *sqlx.DB) *AuthPostgresRepo {
	return &AuthPostgresRepo{db: db}
}

func (r *AuthPostgresRepo) CreateUser(user models.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return fmt.Errorf("походу не сам бади: %s", err.Error())
	}
	// User Create
	var userId int
	var userQuery string = fmt.Sprintf(
		`INSERT INTO %s (email, username, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		usersTable)
	result := tx.QueryRow(userQuery, user.Email, user.Username, user.Password)
	if err := result.Scan(&userId); err != nil {
		_ = tx.Rollback()
		logrus.Errorf("result.Scan(&userId): %s", err.Error())
		return fmt.Errorf("username or email already exists")
	}
	// Wallet Create
	return createWallet(tx, userId)
}

func createWallet(tx *sqlx.Tx, userId int) error {
	var walletId int
	var walletQuery string = fmt.Sprintf(
		`INSERT INTO %s (usd, eur, rub) VALUES (default, default, default) RETURNING id`,
		walletsTable)
	result := tx.QueryRow(walletQuery)
	if err := result.Scan(&walletId); err != nil {
		_ = tx.Rollback()
		logrus.Errorf("%s", err.Error())
		return fmt.Errorf("походу не сам бади: %s", err.Error())
	}
	// Create relationship user<->wallet
	relationQuery := fmt.Sprintf(`INSERT INTO %s (user_id, wallet_id) VALUES ($1, $2)`,
		usersWalletsTable)
	_, err := tx.Exec(relationQuery, userId, walletId)
	if err != nil {
		_ = tx.Rollback()
		logrus.Errorf("tx.Exec(relationQuery, userId, walletId): %s", err.Error())
		return fmt.Errorf("походу не сам бади: %s", err.Error())
	}

	return tx.Commit()
}

func (r *AuthPostgresRepo) GetUser(userInput models.UserLogin) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, userInput.Username, userInput.Password)

	return user, err
}
