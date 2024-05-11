package middleware

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func SetGlobalTimeout(timeout uint8) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path != "/api/v1/sse/handle-sse-without-channel" {
				ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(timeout)*time.Second)
				defer cancel()
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}
