package exception

import (
	"context"
	"encoding/json"
	"net/http"
	modelresponse "project-user/models/responses"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, requestId string, err interface{}) error {
	var httpStatusCode int
	var errorMessage interface{}
	if exception, ok := err.(BadRequestException); ok {
		httpStatusCode = exception.Code
		errorMessage = exception.Error()
	} else if exception, ok := err.(NotFoundException); ok {
		httpStatusCode = exception.Code
		errorMessage = exception.Error()
	} else if exception, ok := err.(ValidationException); ok {
		httpStatusCode = exception.Code
		var exceptionError map[string]interface{}
		json.Unmarshal([]byte(exception.Error()), &exceptionError)
		errorMessage = exceptionError
	} else if exception, ok := err.(TimeoutCancelException); ok {
		httpStatusCode = exception.Code
		errorMessage = exception.Error()
	} else if exception, ok := err.(InternalServerErrorException); ok {
		httpStatusCode = exception.Code
		errorMessage = exception.Error()
	} else {
		httpStatusCode = http.StatusInternalServerError
		errorMessage = "internal server error"
	}
	return modelresponse.ToResponse(c, httpStatusCode, requestId, "", errorMessage)
}

func CheckError(err error) error {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return NewTimeoutCancelException()
	}
	return NewInternalServerErrorException()

}
