package entity

type User struct {
	Login   string
	HashPWD string
	Salt    int
}
