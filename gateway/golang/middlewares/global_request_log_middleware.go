package middleware

import (
	"bytes"
	"fmt"
	helper "gateway/helpers"
	modelresponse "gateway/models/responses"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func SetGlobalRequestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		datetimeNowRequest := time.Now()
		requestMethod := c.Request().Method
		requestId := c.Request().Context().Value(RequestIdKey).(string)

		var rBody string
		rBody = `""`
		if c.Request().Body != nil {
			requestBody, err := io.ReadAll(c.Request().Body)
			defer c.Request().Body.Close()
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			}
			body := io.NopCloser(bytes.NewBuffer(requestBody))
			c.Request().Body = body

			rBody = strings.ReplaceAll(string(requestBody), "\n", "")
			rBody = strings.ReplaceAll(rBody, "\t", "")
		}

		// var tokenAuthorization string
		// cookie, err := c.Request().Cookie("Authorization")
		// if err != nil && !errors.Is(err, http.ErrNoCookie) {
		// 	helper.PrintLogToTerminal(err, requestId)
		// 	response := modelresponse.Response{
		// 		Data:  "",
		// 		Error: "internal server error",
		// 	}
		// 	return c.JSON(http.StatusInternalServerError, response)
		// } else if err != nil && errors.Is(err, http.ErrNoCookie) {
		// 	tokenAuthorization = ""
		// } else {
		// 	tokenAuthorization = cookie.Value
		// }

		host := c.Request().Host
		protocol := ""
		if c.Request().TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}
		urlPath := c.Request().URL.Path
		userAgent := c.Request().Header.Get("User-Agent")
		remoteAddr := c.Request().RemoteAddr
		forwardedFor := c.Request().Header.Get("X-Forwarded-For")

		requestLog := `{"requestTime": "` + datetimeNowRequest.String() + `", "app": "project-gateway", "method": "` + requestMethod + `","requestId":"` + requestId + `","host": "` + host + `","urlPath":"` + urlPath + `","protocol":"` + protocol + `","body": ` + rBody + `, "userAgent": "` + userAgent + `", "remoteAddr": "` + remoteAddr + `", "forwardedFor": "` + forwardedFor + `"}`
		fmt.Println(requestLog)
		return next(c)
	}
}
