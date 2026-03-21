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

	workerService employee.IEmployeeServiceModifier
}

func NewEmployyRouter(
	log *slog.Logger,
	session auth.SessionConfig,
	employeeService employee.IEmployeeServiceModifier,
) EmployeeRouter {
	return EmployeeRouter{
		log:           log,
		Session:       session,
		workerService: employeeService,
	}
}

func (e EmployeeRouter) HandlerCreateEmployee(ctx echo.Context) error {
	return nil
}
