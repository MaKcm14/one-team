package divisions

import (
	"fmt"
	"strconv"
	"strings"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
	"github.com/labstack/echo/v4"
)

const (
	divisionIDQueryParamKey = "division_id"
)

func validateDivisionID(ctx echo.Context) (int, error) {
	val := ctx.QueryParam(divisionIDQueryParamKey)
	if len(val) == 0 {
		return 0, fmt.Errorf("parameter '%s' can't be empty", divisionIDQueryParamKey)
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("parameter '%s' has wrong format", divisionIDQueryParamKey)
	}

	if id <= 0 {
		return 0, fmt.Errorf("paramter '%s' can't be less than 1", divisionIDQueryParamKey)
	}
	return id, nil
}

func validateFilters(ctx echo.Context, pageNum int) (division.Filters, error) {
	var filters division.Filters

	val := ctx.QueryParam(division.NameFilterName)
	if len(val) == 0 {
		return division.Filters{}, fmt.Errorf("parameter '%s' can't be empty", division.NameFilterName)
	}

	parts := strings.Split(val, ":")
	if len(parts) != 2 {
		return division.Filters{}, fmt.Errorf("parameter '%s' has wrong format: it must be set as 'LIKE_NAME:LIKE_TYPE' ",
			division.NameFilterName)
	}

	if !entity.IsDivisionTypeValid(entity.DivisionType(parts[1])) && parts[1] != "" {
		return division.Filters{}, fmt.Errorf("wrong division type was got")
	}

	filters.Names = division.NameFilter{
		Name:    parts[0],
		Type:    entity.DivisionType(parts[1]),
		PageNum: pageNum,
	}
	return filters, nil
}
