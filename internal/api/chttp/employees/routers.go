package employees

import (
	"log/slog"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
	"github.com/labstack/echo/v4"
)

type EmployeeRouter struct {
	log     *slog.Logger
	Session auth.SessionConfig

	workerService employee.IEmployeeService
}

func NewEmployeeRouter(
	log *slog.Logger,
	session auth.SessionConfig,
	employeeService employee.IEmployeeService,
) EmployeeRouter {
	return EmployeeRouter{
		log:           log,
		Session:       session,
		workerService: employeeService,
	}
}

func (e EmployeeRouter) HandlerCreateEmployee(eCtx echo.Context) error {
	return nil
}

func (e EmployeeRouter) HandlerGetTitles(eCtx echo.Context) error {
	return nil
}

func (e EmployeeRouter) HandlerGetCitizenships(eCtx echo.Context) error {
	return nil
}
