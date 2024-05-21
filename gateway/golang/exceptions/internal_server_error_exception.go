package exception

import "net/http"

type InternalServerErrorException struct {
	Code    int
	Message string
}

func NewInternalServerErrorException() InternalServerErrorException {
	return InternalServerErrorException{Code: http.StatusInternalServerError, Message: "internal server error"}
}

func (exception InternalServerErrorException) Error() string {
	return exception.Message
}
