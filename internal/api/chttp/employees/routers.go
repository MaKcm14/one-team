package employees

import (
	"context"
	"errors"
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
	var worker entity.Employee
	if err := eCtx.Bind(&worker); err != nil {
		e.log.Warn(fmt.Sprintf("Warn of parse the request's body: %s", err))
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err := e.workerService.CreateEmployee(ctx, worker)
	if err != nil {
		if errors.Is(err, employee.ErrTitleNotFound) {
			e.log.Warn(fmt.Sprintf("Warn of unknown title in request: %s", err))

			return eCtx.JSON(http.StatusNotFound, server.ErrorResponse{
				Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, ErrUnknownTitle),
			})
		} else if errors.Is(err, employee.ErrCitizenshipNotFound) {
			e.log.Warn(fmt.Sprintf("Warn of unknown citizenship in request: %s", err))

			return eCtx.JSON(http.StatusNotFound, server.ErrorResponse{
				Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, ErrUnknownCitizenship),
			})
		} else if errors.Is(err, employee.ErrEmployeeExists) {
			e.log.Warn(fmt.Sprintf("Warn of creating an existing employee: %s", err))

			return eCtx.JSON(http.StatusConflict, server.ErrorResponse{
				Error: server.ErrAlreadyExists.Error(),
			})
		}
		e.log.Error(fmt.Sprintf("Error of handling the request: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusCreated)
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

func (e EmployeeRouter) HandlerUpdateEmployee(eCtx echo.Context) error {
	var worker entity.Employee
	if err := eCtx.Bind(&worker); err != nil {
		e.log.Error(fmt.Sprintf("Error of binding the request's body: %s", err))
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err := e.workerService.UpdateEmployee(ctx, worker)
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of updating the employee: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusOK)
}
