package exception

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcErrorHandler(err interface{}) error {
	if exception, ok := err.(BadRequestException); ok {
		return status.Error(codes.InvalidArgument, exception.Error())
	} else if exception, ok := err.(NotFoundException); ok {
		return status.Error(codes.NotFound, exception.Error())
	} else if exception, ok := err.(ValidationException); ok {
		return status.Error(codes.InvalidArgument, exception.Error())
	} else if exception, ok := err.(TimeoutCancelException); ok {
		return status.Error(codes.DeadlineExceeded, exception.Error())
	} else if exception, ok := err.(InternalServerErrorException); ok {
		return status.Error(codes.Internal, exception.Error())
	} else {
		return status.Error(codes.Internal, "internal server error")
	}
}
