package exception

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"project-user/helpers"
)

type ResponseException struct {
	Code    int
	Message string
}

func NewResponseException(code int, message string) ResponseException {
	return ResponseException{Code: code, Message: message}
}

func (exception ResponseException) Error() string {
	return exception.Message
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func toErrorMessage(field string, message string) ([]byte, error) {
	var errorMessages []ErrorMessage
	var errorMessage ErrorMessage
	errorMessage.Field = field
	errorMessage.Message = message
	errorMessages = append(errorMessages, errorMessage)
	return json.Marshal(errorMessages)
}

func toErrorMessageRequestTimeout(requestId string) error {
	errorMessagesByte, err := toErrorMessage("message", "time out or user cancel the request")
	if err != nil {
		helpers.PrintLogToTerminal(err, requestId)
		return toErrorMessageInternalServerError(requestId)
	}
	return NewResponseException(http.StatusRequestTimeout, string(errorMessagesByte))
}

func toErrorMessageInternalServerError(requestId string) error {
	errorMessagesByte, err := toErrorMessage("message", "internal server error")
	if err != nil {
		helpers.PrintLogToTerminal(err, requestId)
	}
	return NewResponseException(http.StatusInternalServerError, string(errorMessagesByte))
}

func CheckError(err error, requestId string) error {
	helpers.PrintLogToTerminal(err, requestId)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return toErrorMessageRequestTimeout(requestId)
	} else {
		return toErrorMessageInternalServerError(requestId)
	}
}

func ToResponseExceptionRequestValidation(requestId string, validationErrorMessages []helpers.ErrorMessage) error {
	var validationErrorMessageByte []byte
	validationErrorMessageByte, err := json.Marshal(validationErrorMessages)
	if err != nil {
		return CheckError(err, requestId)
	}
	err = errors.New(string(validationErrorMessageByte))
	helpers.PrintLogToTerminal(err, requestId)
	return NewResponseException(http.StatusBadRequest, string(validationErrorMessageByte))
}

func ToResponseException(err error, requestId string, httpCode int, message string) error {
	helpers.PrintLogToTerminal(err, requestId)
	field := "message"
	errorMessagesByte, err := toErrorMessage(field, message)
	if err != nil {
		helpers.PrintLogToTerminal(err, requestId)
		return CheckError(err, requestId)
	}
	return NewResponseException(httpCode, string(errorMessagesByte))
}
