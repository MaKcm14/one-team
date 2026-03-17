package mw

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetLogMsg(key string, val any) string {
	return fmt.Sprintf("%s:%s\n", key, val)
}

func formatHeaders(headers http.Header) []string {
	res := make([]string, 0, len(headers))
	for key, val := range headers {
		res = append(res, GetLogMsg(key, val))
	}
	return res
}

func formatBody(r io.Reader) (formatStr string, body []byte) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "None", nil
	}
	msg := string(body)

	return msg, body
}

func extractRequestInfo(ctx echo.Context) string {
	b := strings.Builder{}

	b.WriteString(GetLogMsg("IP", ctx.Request().RemoteAddr))
	b.WriteString(GetLogMsg("METHOD", ctx.Request().Method))
	b.WriteString(GetLogMsg("URI", ctx.Request().RequestURI))
	b.WriteString(GetLogMsg("HEADERS", formatHeaders(ctx.Request().Header)))

	val, body := formatBody(ctx.Request().Body)
	ctx.Request().Body.Close()

	if body != nil {
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))
	}
	b.WriteString(GetLogMsg("BODY", val))

	return b.String()
}

func extractResponseInfo(ctx echo.Context) string {
	b := strings.Builder{}

	b.WriteString(GetLogMsg("IP", ctx.Request().RemoteAddr))
	b.WriteString(GetLogMsg("URI", ctx.Request().RequestURI))
	b.WriteString(GetLogMsg("STATUS", ctx.Response().Status))
	b.WriteString(GetLogMsg("HEADERS", formatHeaders(ctx.Response().Header())))
	return b.String()
}
