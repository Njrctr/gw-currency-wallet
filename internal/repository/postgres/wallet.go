package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Njrctr/gw-currency-wallet/internal/models"
	"github.com/jmoiron/sqlx"
)

var ErrBalanceCheck = "pq: Недостаточно средст для выполнения операции"

type WalletPostgresRepo struct {
	db *sqlx.DB
}

func NewWalletPostgresRepo(db *sqlx.DB) *WalletPostgresRepo {
	return &WalletPostgresRepo{db: db}
}

func (r *WalletPostgresRepo) GetWallet(ctx context.Context, userId int) (models.Wallet, error) {
	var wallet models.Balance
	query := fmt.Sprintf(
		"SELECT usd, rub, eur FROM %s w INNER JOIN %s uw on w.id=uw.wallet_id WHERE uw.user_id = $1",
		walletsTable, usersWalletsTable)
	err := r.db.Get(&wallet, query, userId)

	return models.Wallet{Balance: wallet}, err
}

func (r *WalletPostgresRepo) EditBalance(
	ctx context.Context,
	userId int, operation models.EditWallet,
) (models.Balance, error) {

	var typeQuery string
	if operation.OperationType == "DEPOSIT" {
		typeQuery = fmt.Sprintf("%s=%s+%v", operation.Currency, operation.Currency, operation.Amount)
	} else {
		typeQuery = fmt.Sprintf("%s=%s-%v", operation.Currency, operation.Currency, operation.Amount)
	}

	var newBalance models.Balance
	query := fmt.Sprintf(
		`UPDATE %s w 
		SET %s 
		FROM %s uw 
		WHERE w.id=uw.wallet_id 
		AND uw.user_id=$1 
		RETURNING w.usd, w.rub, w.eur`, walletsTable, typeQuery, usersWalletsTable)

	err := r.db.Get(&newBalance, query, userId)
	if err != nil && err.Error() == ErrBalanceCheck {
		return models.Balance{}, errors.New("недостаточно средств на счете")
	}

	return newBalance, err
}

func (r *WalletPostgresRepo) Transfer(ctx context.Context, input models.TransferOperation) (models.Balance, error) {
	var newBalance models.Balance
	transer := fmt.Sprintf("%s=%s-%v, %s=%s+%v", input.From, input.From, input.Amount,
		input.To, input.To, input.Amount*input.Rate)
	query := fmt.Sprintf(
		`UPDATE %s w 
		SET %s 
		FROM %s uw
		WHERE w.id=uw.wallet_id
		AND uw.user_id=$1 
		RETURNING w.%s, w.%s`, walletsTable, transer, usersWalletsTable, input.From, input.To)

	fmt.Println(query)
	err := r.db.Get(&newBalance, query, input.UserId)
	if err != nil && err.Error() == ErrBalanceCheck {
		return models.Balance{}, errors.New("недостаточно средств на счете")
	}

	return newBalance, nil
}
