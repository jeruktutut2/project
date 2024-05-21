package exception

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
