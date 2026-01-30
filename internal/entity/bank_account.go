package entity

type BankAccountID int

type BankAccount struct {
	ID    BankAccountID `json:"id"`
	Money float64       `json:"money"`
}
