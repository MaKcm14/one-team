package employees

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
	"github.com/labstack/echo/v4"
)

const (
	// QueryParamKeys.
	pageNumQueryParamKey          = "page_num"
	passportDataQueryParamKey     = "passport_data"
	salaryDownBoundQueryParamKey  = "down"
	salaryUpperBoundQueryParamKey = "up"
	titleIDQueryParamKey          = "title_id"
	filtersQueryParamKey          = "filters"
)

func validateSalaryDownBound(ctx echo.Context) (int, error) {
	param := ctx.QueryParam(salaryDownBoundQueryParamKey)
	if len(param) == 0 {
		return 0, fmt.Errorf("error of '%s' query param value: wasn't set", salaryDownBoundQueryParamKey)
	}

	downBound, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("error of converting the '%s' query param into num", salaryDownBoundQueryParamKey)
	} else if downBound < 0 {
		return 0, fmt.Errorf("error of '%s' query param: can't be less than 0", salaryDownBoundQueryParamKey)
	}
	return downBound, nil
}

func validateSalaryUpperBound(ctx echo.Context) (int, error) {
	param := ctx.QueryParam(salaryUpperBoundQueryParamKey)
	if len(param) == 0 {
		return 0, fmt.Errorf("error of '%s' query param value: wasn't set", salaryUpperBoundQueryParamKey)
	}

	upperBound, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("error of converting the '%s' query param into num", salaryUpperBoundQueryParamKey)
	} else if upperBound < 0 {
		return 0, fmt.Errorf("error of '%s' query param: can't be less than 0", salaryDownBoundQueryParamKey)
	}
	return upperBound, nil
}

func validateTitleID(ctx echo.Context) (int, error) {
	param := ctx.QueryParam(titleIDQueryParamKey)
	if len(param) == 0 {
		return 0, fmt.Errorf("error of '%s' query param value: wasn't set", titleIDQueryParamKey)
	}

	titleID, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("error of converting the '%s' query param into ID", titleIDQueryParamKey)
	} else if titleID < 0 {
		return 0, fmt.Errorf("error of '%s' query param: can't be less than 0", titleIDQueryParamKey)
	}
	return titleID, nil
}

func validatePassportData(data string) error {
	rule, err := regexp.Compile(`^\d{4} \d{6}$`)
	if err != nil {
		return fmt.Errorf("error of template compile")
	}

	if !rule.Match([]byte(data)) {
		return fmt.Errorf("passport data doesn't have valid format")
	}
	return nil
}

func validateFilters(ctx echo.Context, pageNum int) (employee.Filter, error) {
	filter := employee.Filter{}

	val := ctx.QueryParam(employee.NamesFilterName)
	names := strings.Split(val, ":")
	if len(names) == 3 {
		filter.Names.IsActive = true
		filter.Names.FirstName = names[0]
		filter.Names.LastName = names[1]
		filter.Names.Patronymic = names[2]
	}

	val = ctx.QueryParam(employee.PassportDataFilterName)
	if len(val) != 0 {
		filter.Passport.IsActive = true
		filter.Passport.PassportData = val
	}

	val = ctx.QueryParam(employee.UnitFilterName)
	if len(val) != 0 {
		filter.Unit.IsActive = true
		filter.Unit.Name = val

		val = ctx.QueryParam(employee.UnitTypeFilterName)
		if entity.IsDivisionTypeValid(entity.DivisionType(val)) {
			return employee.Filter{}, fmt.Errorf("'%s' is wrong: can't recognize it", employee.UnitTypeFilterName)
		}
		filter.Unit.Type = entity.DivisionType(val)
	}

	if !filter.Names.IsActive &&
		!filter.Passport.IsActive && !filter.Unit.IsActive {
		return employee.Filter{}, fmt.Errorf("unknown filter format")
	}

	filter.Names.PageNum = pageNum
	filter.Passport.PageNum = pageNum
	filter.Unit.PageNum = pageNum

	return filter, nil
}

func validatePageNum(ctx echo.Context) (int, error) {
	val := ctx.QueryParam(pageNumQueryParamKey)
	if len(val) == 0 {
		return 0, fmt.Errorf("wrong '%s' value: can't be empty", pageNumQueryParamKey)
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("wrong format for '%s' parameter: can't parse it", pageNumQueryParamKey)
	}

	if num < 0 {
		return 0, fmt.Errorf("wrong '%s' value: can't be less than 0", pageNumQueryParamKey)
	}
	return num, nil
}
