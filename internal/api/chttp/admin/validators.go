package admin

import (
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
	"github.com/labstack/echo/v4"
)

func validateFilters(ctx echo.Context) (user.Filters, error) {
	pageNum, err := server.ValidatePageNum(ctx)
	if err != nil {
		return user.Filters{}, err
	}

	login, err := server.ValidateLoginQueryParam(ctx)
	if err != nil {
		return user.Filters{}, err
	}

	return user.Filters{
		LoginFilter: user.LoginFilter{
			IsActive: true,
			PageNum:  pageNum,
			Login:    login,
		},
	}, nil
}
