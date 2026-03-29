package server

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	pageNumQueryParamKey = "page_num"
)

func ValidatePageNum(ctx echo.Context) (int, error) {
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
