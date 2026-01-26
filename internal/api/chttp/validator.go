package chttp

import (
	"regexp"
	"strconv"

	"github.com/labstack/echo"

	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

const (
	userIDParamName = "userID"
)

func validateUserID(ctx echo.Context) (entity.UserID, error) {
	id, err := strconv.Atoi(ctx.QueryParam(userIDParamName))
	if err != nil {
		return 0, ErrRequestData
	}
	return entity.UserID(id), nil
}

func validateUser(userCfg repo.UserConfig) error {
	passReg := regexp.MustCompile(`^\d{4} \d{6}$`)
	if !passReg.MatchString(userCfg.Passport) {
		return ErrRequestData
	}
	return nil
}
