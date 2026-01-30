package entity

type UserID int

type UserProfile struct {
	AdminStatus bool   `json:"admin"`
	PwdHash     []byte `json:"-"`
}

type User struct {
	ID       UserID      `json:"id"`
	Name     string      `json:"name"`
	Surname  string      `json:"surname"`
	Passport string      `json:"passport"`
	Account  BankAccount `json:"bank_account"`
	Profile  UserProfile `json:"profile"`
}
