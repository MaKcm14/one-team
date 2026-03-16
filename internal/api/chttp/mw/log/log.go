package log

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/mw"
)

func LoggerMW(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			log.Info(fmt.Sprintf("REQUEST: %s\n", mw.ExtractRequestInfo(ctx)))

			err := next(ctx)

			msg := fmt.Sprintf("RESPONSE: %s\n", mw.ExtractResponseInfo(ctx))
			if err != nil {
				msg += mw.GetLogMsg("ERR", err)
			}
			log.Info(msg)

			return nil
		}
	}
}
