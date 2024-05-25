package interceptor

import (
	"context"
	"encoding/json"
	"errors"
	exception "project-user/exceptions"
	helper "project-user/helpers"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SetLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("cannot find metadata")
		helper.PrintLogToTerminal(err, "requestid")
		return nil, exception.GrpcErrorHandler(err)
	}

	requestIdA, ok := md["requestid"]
	if !ok {
		err = errors.New("cannot find requestid")
		helper.PrintLogToTerminal(err, "requestid")
		return nil, exception.GrpcErrorHandler(err)
	}

	sessionidA, ok := md["sessionid"]
	if !ok {
		err = errors.New("cannot find sessionid")
		helper.PrintLogToTerminal(err, requestIdA[0])
		return nil, exception.GrpcErrorHandler(requestIdA)
	}

	var reqString string
	reqString = ``
	if info.FullMethod != "/protofiles.UserService/Login" {
		reqByte, err := json.Marshal(req)
		if err != nil {
			helper.PrintLogToTerminal(err, requestIdA[0])
			return nil, exception.GrpcErrorHandler(err)
		}
		reqString = string(reqByte)
	}

	requestLog := `{"grpcRequestTime": "` + time.Now().String() + `", "app": "project-user" ,"requestId": "` + requestIdA[0] + `", "sessionid": "` + sessionidA[0] + `", "urlPath":"` + info.FullMethod + `","body": "` + reqString + `", "metadata:": ` + md.String() + `}`
	println(requestLog)

	h, errHandler := handler(ctx, req)
	hByte, err := json.Marshal(h)
	if err != nil {
		helper.PrintLogToTerminal(err, requestIdA[0])
		return nil, exception.GrpcErrorHandler(err)
	}

	var hString string
	hString = `""`
	if hByte == nil {
		hString = `""`
	} else {
		hString = string(hByte)
		if hString == "null" {
			hString = `""`
		}
	}

	responseLog := `{"grpcResponseTime": "` + time.Now().String() + `", "app": "project-user", "requestId": "` + requestIdA[0] + `", "body": ` + hString + `}`
	println(responseLog)

	return h, errHandler
}
