package exception

import "net/http"

type NotFoundException struct {
	Code    int
	Message string
}

func NewNotFoundException(message string) NotFoundException {
	return NotFoundException{Code: http.StatusNotFound, Message: message}
}

func (exception NotFoundException) Error() string {
	return exception.Message
}
