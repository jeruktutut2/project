package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	helper "project-user/helpers"
	modelresponse "project-user/models/responses"
	"time"

	"github.com/labstack/echo/v4"
)

func SetGlobalRequestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		datetimeNowRequest := time.Now()
		requestMethod := c.Request().Method
		requestId := c.Request().Context().Value(RequestIdKey).(string)

		var requestBody string
		requestBody = `""`
		if c.Request().GetBody != nil {
			jsonRequestBodyMap := make(map[string]interface{})
			err := json.NewDecoder(c.Request().Body).Decode(&jsonRequestBodyMap)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			}
			jsonRequestBodyByte, err := json.Marshal(jsonRequestBodyMap)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "cannot convert request body to json")
			}
			requestBody = string(jsonRequestBodyByte)
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

		requestLog := `{"requestTime": "` + datetimeNowRequest.String() + `", "method": "` + requestMethod + `","requestId":"` + requestId + `","host": "` + host + `","urlPath":"` + urlPath + `","protocol":"` + protocol + `","body": ` + requestBody + `, "userAgent": "` + userAgent + `", "remoteAddr": "` + remoteAddr + `", "forwardedFor": "` + forwardedFor + `"}`
		fmt.Println(requestLog)
		return next(c)
	}
}
