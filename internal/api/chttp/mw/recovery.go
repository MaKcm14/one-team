package mw

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
)

func Recovery(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					b := strings.Builder{}

					b.WriteString(fmt.Sprintf("RECOVER: catch error: %s\n", err))
					b.WriteString(fmt.Sprintf("REQUEST: %s\n", ExtractRequestInfo(ctx)))
					b.WriteString(fmt.Sprintf("RESPONSE: %s\n", ExtractResponseInfo(ctx)))

					log.Error(b.String())
				}
			}()
			return next(ctx)
		}
	}
}
