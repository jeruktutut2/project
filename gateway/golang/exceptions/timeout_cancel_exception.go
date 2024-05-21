package exception

import "net/http"

type TimeoutCancelException struct {
	Code    int
	Message string
}

func NewTimeoutCancelException() TimeoutCancelException {
	return TimeoutCancelException{Code: http.StatusRequestTimeout, Message: "time out or user cancel the request"}
}

func (exception TimeoutCancelException) Error() string {
	return exception.Message
}
