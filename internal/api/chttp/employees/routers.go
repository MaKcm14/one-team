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

	err := validatePassportData(worker.PassportData)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err = e.workerService.CreateEmployee(ctx, worker)
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

func (e EmployeeRouter) HandlerCountEmployeeWithCitizenship(eCtx echo.Context) error {
	type response struct {
		Statistics []employee.EmployeeCitizenshipStatistic `json:"statistics"`
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	stats, err := e.workerService.CountEmployeesWithCitizenship(ctx)
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of counting the stats: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Statistics: stats,
	})
}

func (e EmployeeRouter) HandlerCountEmployeesWithSalaryBoundary(eCtx echo.Context) error {
	type response struct {
		Count int `json:"count"`
	}

	downBound, err := validateSalaryDownBound(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	upperBound, err := validateSalaryUpperBound(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	if downBound > upperBound {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	titleID, err := validateTitleID(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	count, err := e.workerService.CountEmployeesWithSalaryBounds(
		ctx,
		titleID,
		employee.SalaryBounds{
			DownBoundary: downBound,
			UpBoundary:   upperBound,
		})
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of counting the statistics: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Count: count,
	})
}

func (e EmployeeRouter) HandlerGetEmployeeWithFilter(eCtx echo.Context) error {
	type response struct {
		Employees []entity.Employee `json:"employees"`
	}

	pageNum, err := validatePageNum(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	filters, err := validateFilters(eCtx, pageNum)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	employeeList, err := e.workerService.GetEmployeesWithFilters(ctx, filters, pageNum)
	if err != nil {
		e.log.Error(fmt.Sprintf("Error of getting the employees: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Employees: employeeList,
	})
}

func (e EmployeeRouter) HandlerDeleteEmployee(eCtx echo.Context) error {
	const reportName = "deleted_employee_report.xlsx"

	employeeID, err := validateEmployeeID(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	path, err := e.workerService.DeleteEmployee(ctx, employeeID)
	if err != nil {
		if errors.Is(err, employee.ErrReportCreating) {
			e.log.Error(fmt.Sprintf("Error of creating the report for deleted employee: %s", err))
		} else {
			e.log.Error(fmt.Sprintf("Error of deleting the employee: %s", err))
		}
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	eCtx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	eCtx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", reportName))

	return eCtx.Attachment(path, reportName)
}
