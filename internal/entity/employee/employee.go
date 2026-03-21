package entity

import (
	"time"

	division "github.com/MaKcm14/one-team/internal/entity/division"
)

type Title string

type Citizenship string

type Employee struct {
	TinNum       string
	Snils        string
	PassportData string
	PhoneNum     string
	FirstName    string
	LastName     string
	Patronymic   string
	Address      string
	Title        Title
	Citizenship  Citizenship
	HiringDate   time.Time
	UnitID       division.DivisionID
	Education    string
	Salary       float64
}
