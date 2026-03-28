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
