package divisions

import (
	"context"
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
