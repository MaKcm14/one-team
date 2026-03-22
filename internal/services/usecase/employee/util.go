package employee

import entity "github.com/MaKcm14/one-team/internal/entity/employee"

type EmployeeCitizenshipStatistic struct {
	Citizenship   entity.Citizenship `json:"citizenship"`
	EmployeeCount int                `json:"employee_count"`
}
