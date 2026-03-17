package mw

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func LoggerMW(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			log.Info(fmt.Sprintf("REQUEST: %s\n", extractRequestInfo(ctx)))

			err := next(ctx)

			msg := fmt.Sprintf("RESPONSE: %s\n", extractResponseInfo(ctx))
			if err != nil {
				msg += GetLogMsg("ERR", err)
			}
			log.Info(msg)

			return nil
		}
	}
}
