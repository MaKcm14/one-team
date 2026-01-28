package entity

type UserID int

type User struct {
	ID       UserID      `json:"id"`
	Name     string      `json:"name"`
	Surname  string      `json:"surname"`
	Passport string      `json:"passport"`
	Account  BankAccount `json:"bank_account"`
	PwdHash  []byte      `json:"-"`
}
