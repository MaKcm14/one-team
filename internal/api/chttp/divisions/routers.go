package divisions

import (
	"log/slog"

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
	return nil
}
