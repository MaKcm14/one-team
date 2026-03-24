package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type IEmployeeServiceWriter interface {
	CreateEmployee(ctx context.Context, employee entity.Employee) error
	UpdateEmployee(ctx context.Context, employee entity.Employee) error
}

type IEmployeeServiceReader interface {
	CountEmployeesWithCitizenship(ctx context.Context) ([]EmployeeCitizenshipStatistic, error)
	CountEmployeesWithSalaryBounds(ctx context.Context, titleID int, bounds SalaryBounds) (int, error)
	GetEmployeesWithFilters(ctx context.Context, filters Filter, pageNum int) ([]entity.Employee, error)
}

type IEmployeeTitleReader interface {
	GetTitles(ctx context.Context) ([]entity.Title, error)
}

type IEmployeeCitizenshipReader interface {
	GetCitizenships(ctx context.Context) ([]entity.Citizenship, error)
}

type IEmployeeService interface {
	IEmployeeServiceWriter
	IEmployeeServiceReader

	IEmployeeCitizenshipReader
	IEmployeeTitleReader
}
