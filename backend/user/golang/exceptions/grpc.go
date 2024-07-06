package exception

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcErrorHandler(err interface{}) error {
	if exception, ok := err.(ResponseException); ok {
		if exception.Code == http.StatusBadRequest {
			return status.Error(codes.InvalidArgument, exception.Error())
		} else if exception.Code == http.StatusNotFound {
			return status.Error(codes.NotFound, exception.Error())
		} else if exception.Code == http.StatusRequestTimeout {
			return status.Error(codes.DeadlineExceeded, exception.Error())
		} else if exception.Code == http.StatusInternalServerError {
			return status.Error(codes.Internal, exception.Error())
		} else {
			return status.Error(codes.Internal, "internal server error")
		}
	} else {
		return status.Error(codes.Internal, "internal server error")
	}
}
