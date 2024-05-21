package exception

import "net/http"

type BadRequestException struct {
	Code    int
	Message string
}

func NewBadRequestException(message string) BadRequestException {
	return BadRequestException{Code: http.StatusBadRequest, Message: message}
}

func (exception BadRequestException) Error() string {
	return exception.Message
}
