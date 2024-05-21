package modelresponse

import (
	"encoding/json"
	"fmt"
	helper "gateway/helpers"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func ToResponse(c echo.Context, httpStatusCode int, requestId string, responseData interface{}, responseError interface{}) error {
	r := Response{
		Data:  responseData,
		Error: responseError,
	}
	respByte, err := json.Marshal(r)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		response := `{"data": "", "error": "internal server error"}`
		return c.JSON(http.StatusInternalServerError, response)
	}
	responseBody := string(respByte)
	log := `{"responseTime": "` + time.Now().String() + `", "app": "project-gateway", "requestId": "` + requestId + `", "response": ` + responseBody + `}`
	fmt.Println(log)
	return c.JSON(httpStatusCode, r)
}
