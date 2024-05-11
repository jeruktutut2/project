package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetGlobalRequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		randomUUID := uuid.New().String()
		ctx := context.WithValue(c.Request().Context(), RequestIdKey, randomUUID)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
