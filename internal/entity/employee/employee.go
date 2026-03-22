package entity

import (
	"time"

	division "github.com/MaKcm14/one-team/internal/entity/division"
)

type Title struct {
	ID   int    `json:"title_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Citizenship struct {
	ID   int    `json:"citizenship_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Employee struct {
	EmployeeID   int               `json:"employee_id,omitempty"`
	TinNum       string            `json:"tin_num,omitempty"`
	Snils        string            `json:"snils,omitempty"`
	PassportData string            `json:"passport_data,omitempty"`
	PhoneNum     string            `json:"phone_num,omitempty"`
	FirstName    string            `json:"first_name,omitempty"`
	LastName     string            `json:"last_name,omitempty"`
	Patronymic   string            `json:"patronymic,omitempty"`
	Address      string            `json:"address,omitempty"`
	Title        Title             `json:"title,omitempty"`
	Citizenship  Citizenship       `json:"citizenship,omitempty"`
	HiringDate   time.Time         `json:"hiring_date,omitempty"`
	Unit         division.Division `json:"unit,omitempty"`
	Education    string            `json:"education,omitempty"`
	Salary       float64           `json:"salary,omitempty"`
}
