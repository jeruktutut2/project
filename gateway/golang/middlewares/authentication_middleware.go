package middleware

import (
	"context"
	helper "gateway/helpers"
	modelresponse "gateway/models/responses"
	util "gateway/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func Authenticate(redisUtil util.RedisUtil) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Context().Value(RequestIdKey).(string)
			cookie, err := c.Cookie("Authorization")
			if err != nil && err != http.ErrNoCookie {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			} else if err != nil && err == http.ErrNoCookie {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized")
			}
			key := cookie.Value
			_, err = redisUtil.GetClient().Get(c.Request().Context(), key).Result()
			if err != nil && err != redis.Nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			} else if err != nil && err == redis.Nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized")
			}
			ctx := context.WithValue(c.Request().Context(), SessionIdKey, key)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func GetSessionId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := c.Request().Context().Value(RequestIdKey).(string)
		var key string
		cookie, err := c.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			helper.PrintLogToTerminal(err, requestId)
			return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
		} else if err != nil && err == http.ErrNoCookie {
			key = ""
		} else {
			key = cookie.Value
		}
		ctx := context.WithValue(c.Request().Context(), SessionIdKey, key)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
