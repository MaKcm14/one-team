package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type IEmployeeRepoReader interface {
	IsEmployeeExists(ctx context.Context, worker entity.Employee) error
	CountEmployeesWithCitizenship(ctx context.Context) ([]EmployeeCitizenshipStatistic, error)
	CountEmployeesWithSalaryBounds(ctx context.Context, bounds SalaryBounds) (int, error)
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
