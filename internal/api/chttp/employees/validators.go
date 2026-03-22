package employees

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	salaryDownBoundQueryParamKey  = "down"
	salaryUpperBoundQueryParamKey = "up"
	titleIDQueryParamKey          = "title_id"
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
