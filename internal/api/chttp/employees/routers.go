package employees

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	entity "github.com/MaKcm14/one-team/internal/entity/employee"
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
	type response struct {
		Titles []entity.Title `json:"titles"`
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	titles, err := e.workerService.GetTitles(ctx)
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of getting the titles: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Titles: titles,
	})
}

func (e EmployeeRouter) HandlerGetCitizenships(eCtx echo.Context) error {
	type response struct {
		Citizenships []entity.Citizenship `json:"citizenships"`
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	citizenships, err := e.workerService.GetCitizenships(ctx)
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of getting the titles: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Citizenships: citizenships,
	})
}
