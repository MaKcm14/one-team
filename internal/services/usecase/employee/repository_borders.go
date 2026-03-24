package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type IEmployeeRepoReader interface {
	IsEmployeeExists(ctx context.Context, worker entity.Employee) error
	CountEmployeesWithCitizenship(ctx context.Context) ([]EmployeeCitizenshipStatistic, error)
	CountEmployeesWithSalaryBounds(ctx context.Context, titleID int, bounds SalaryBounds) (int, error)

	GetEmployeesByName(ctx context.Context, filter NamesFilter) ([]entity.Employee, error)
	GetEmployeesByPassportData(ctx context.Context, filter PassportFilter) ([]entity.Employee, error)

	GetEmployeesByNameInDivision(ctx context.Context, filter NamesFilter, div UnitFilter) ([]entity.Employee, error)
	GetEmployeesByPassportDataInDivision(ctx context.Context, filters PassportFilter, div UnitFilter) ([]entity.Employee, error)
}

type IEmployeeRepoWriter interface {
	CreateEmployee(ctx context.Context, worker entity.Employee) error
	UpdateEmployee(ctx context.Context, worker entity.Employee) error
}

type ITitleRepoReader interface {
	GetTitles(ctx context.Context) ([]entity.Title, error)
}

type ICitizenshipRepoReader interface {
	GetCitizenships(ctx context.Context) ([]entity.Citizenship, error)
}

type IRepoReader interface {
	IEmployeeRepoReader
	ITitleRepoReader
	ICitizenshipRepoReader
}

type IRepoWriter interface {
	IEmployeeRepoWriter
}

type IEmployeeRepo interface {
	IRepoReader
	IRepoWriter
}
