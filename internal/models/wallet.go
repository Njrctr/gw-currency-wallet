package models

type Wallet struct {
	Balance Balance `json:"balance"`
}

type Balance struct {
	Usd float64 `json:"USD" db:"usd"`
	Rub float64 `json:"RUB" db:"rub"`
	Eur float64 `json:"EUR" db:"eur"`
}

type EditWallet struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	OperationType string  `json:"-"` // DEPOSIT or WITHDRAW
}

type Exchange struct {
	From   string  `json:"from_currency"`
	To     string  `json:"to_currency"`
	Amount float64 `json:"amount"`
}

type TransferOperation struct {
	UserId int
	From   string
	To     string
	Amount float64
	Rate   float64
}
