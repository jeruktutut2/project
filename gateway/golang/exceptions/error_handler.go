package exception

import (
	"context"
	modelresponse "gateway/models/responses"
	"net/http"

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
