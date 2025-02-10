package models

type Rates struct {
	Usd float64 `json:"USD" db:"usd"`
	Rub float64 `json:"RUB" db:"rub"`
	Eur float64 `json:"EUR" db:"eur"`
}

type ExchangeRequest struct {
	From   string  `json:"from_currency"`
	To     string  `json:"to_currency"`
	Amount float64 `json:"amount"`
}
