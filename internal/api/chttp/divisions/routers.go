package divisions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
	"github.com/labstack/echo/v4"
)

type DivisionRouter struct {
	log             *slog.Logger
	divisionService division.IDivisionService
}

func NewDivisionRouter(
	log *slog.Logger,
	divisionService division.IDivisionService,
) DivisionRouter {
	return DivisionRouter{
		log:             log,
		divisionService: divisionService,
	}
}

func (d DivisionRouter) HandlerGetDivisions(eCtx echo.Context) error {
	type response struct {
		Divisions []entity.Division `json:"divisions"`
	}

	pageNum, err := server.ValidatePageNum(eCtx)
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

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	divisions, err := d.divisionService.GetDivisions(ctx, filters)
	if err != nil {
		d.log.Error(fmt.Sprintf("Error of getting the titles: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Divisions: divisions,
	})
}

func (d DivisionRouter) HandlerCreateDivision(eCtx echo.Context) error {
	var div entity.Division
	if err := eCtx.Bind(&div); err != nil {
		d.log.Warn(fmt.Sprintf("Warn of binding the request body with division: %s", err))
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	if !entity.IsDivisionTypeValid(div.Type) {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: division type is not valid", server.ErrRequestInfo),
		})
	} else if div.Type == entity.DivisionTypeName && div.SuperdivisionID != 0 {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: superdivision_id is not valid for type 'division'", server.ErrRequestInfo),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err := d.divisionService.CreateDivision(ctx, div)
	if err != nil {
		if errors.Is(err, division.ErrDivisionExists) {
			d.log.Warn("Warn of creating the existing division")
			return eCtx.JSON(http.StatusConflict, server.ErrorResponse{
				Error: fmt.Sprintf("%s: division already exists", server.ErrRequestInfo),
			})
		} else if errors.Is(err, division.ErrDivisionNotFound) {
			d.log.Warn("Warn of creating the division as subdivision of unexisting division")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the superdivision doesn't exist", server.ErrRequestInfo),
			})
		} else if errors.Is(err, division.ErrWrongDivisionsRelation) {
			d.log.Warn("Warn of creating the division with the wrong inter-division types' relation")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the relation between divisions is wrong", server.ErrRequestInfo),
			})
		}

		d.log.Error(fmt.Sprintf("Error of creating the division: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusCreated)
}

func (d DivisionRouter) HandlerDeleteDivision(eCtx echo.Context) error {
	id, err := validateDivisionID(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err = d.divisionService.DeleteDivision(ctx, id)
	if err != nil {
		if errors.Is(err, division.ErrDivisionNotFound) {
			d.log.Warn("Warn of deleting the unexisting division")
			return eCtx.JSON(http.StatusNotFound, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the division doesn't exist", ErrDivisionDeleting),
			})

		} else if errors.Is(err, division.ErrDivisionNotEmpty) {
			d.log.Warn("Warn of deleting the division with employees")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the division has employees", ErrDivisionDeleting),
			})

		} else if errors.Is(err, division.ErrDivisionIsSuperdivision) {
			d.log.Warn("Warn of deleting the division that has sub-divisions")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the division is superdivision for some divisions", ErrDivisionDeleting),
			})
		}

		d.log.Error(fmt.Sprintf("Error of deleting the division: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusOK)
}

func (d DivisionRouter) HandlerUpdateDivision(eCtx echo.Context) error {
	var div entity.Division
	if err := eCtx.Bind(&div); err != nil {
		d.log.Warn("Warn of binding the request body for division updating")
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: the request body has wrong format", server.ErrRequestInfo),
		})
	}

	err := validateDivision(div)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err = d.divisionService.UpdateDivision(ctx, div)
	if err != nil {
		if errors.Is(err, division.ErrDivisionNotFound) {
			d.log.Warn("Warn of updating the unexisting division")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the division doesn't exist", ErrDivisionUpdating),
			})

		} else if errors.Is(err, division.ErrSuperdivisionNotFound) {
			d.log.Warn("Warn of updating the division by unexisting superdivision")
			return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
				Error: fmt.Sprintf("%s: the division of type '%s' may not have the unexisting superdivision",
					ErrDivisionUpdating, div.Type),
			})
		}
		d.log.Error(fmt.Sprintf("Error of updating the division: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}
	return eCtx.NoContent(http.StatusOK)
}

func (d DivisionRouter) HandlerGetSalaryStatisticsOfDivision(eCtx echo.Context) error {
	type response struct {
		SalaryStats division.SalaryStatistics `json:"salary_statistics"`
	}

	id, err := validateDivisionID(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	stats, err := d.divisionService.GetSalaryStatisticsOfDivision(ctx, id)
	if err != nil {
		d.log.Error(fmt.Sprintf("Error of getting the salary statistics of division: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		SalaryStats: stats,
	})
}

func (d DivisionRouter) HandlerGetStateSizeStatisticsOfDivisions(eCtx echo.Context) error {
	type response struct {
		StateSizeStats division.StateSizeStatistics `json:"state_size_statistics"`
	}

	divType, err := validateDivisionType(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	stats, err := d.divisionService.GetStateSizeStatisticsOfDivisions(ctx, divType)
	if err != nil {
		d.log.Error(fmt.Sprintf("Error of getting the state size stats: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		StateSizeStats: stats,
	})
}
