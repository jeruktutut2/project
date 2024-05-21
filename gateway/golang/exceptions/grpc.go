package exception

import (
	"encoding/json"
	modelresponse "gateway/models/responses"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Result struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GrpcErrorHandler(c echo.Context, requestId string, err error) error {
	var httpStatusCode int
	var errorMessage interface{}
	st, ok := status.FromError(err)
	if !ok {
		return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
	}
	errorMessage = st.Message()
	if st.Code() == codes.InvalidArgument {
		httpStatusCode = http.StatusBadRequest
		var exceptionError []Result
		errJsonUnmarshal := json.Unmarshal([]byte(st.Message()), &exceptionError)
		if errJsonUnmarshal != nil {
			errorMessage = st.Message()
		} else {
			errorMessage = exceptionError
		}
	} else if st.Code() == codes.NotFound {
		httpStatusCode = http.StatusNotFound
	} else if st.Code() == codes.DeadlineExceeded {
		httpStatusCode = http.StatusRequestTimeout
	} else if st.Code() == codes.Internal {
		httpStatusCode = http.StatusInternalServerError
	} else {
		httpStatusCode = http.StatusInternalServerError
		errorMessage = "internal server error"
	}
	return modelresponse.ToResponse(c, httpStatusCode, requestId, "", errorMessage)
}
