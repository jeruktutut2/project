package exception

import "net/http"

type ValidationException struct {
	Code    int
	Message string
}

func NewValidationException(message string) ValidationException {
	return ValidationException{Code: http.StatusBadRequest, Message: message}
}

func (exception ValidationException) Error() string {
	return exception.Message
}
