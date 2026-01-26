package entity

type AccountID int

type BankAccount struct {
	ID    AccountID `json:"id"`
	Money float64   `json:"money"`
}
