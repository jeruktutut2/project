package grpcuser

import (
	"context"
	"errors"
	exception "project-user/exceptions"
	pbuser "project-user/grpc/pb/api/v1/user"
	helper "project-user/helpers"
	modelrequest "project-user/models/requests"
	service "project-user/services"

	"google.golang.org/grpc/metadata"
)

type UserGrpcServiceImplementation struct {
	pbuser.UnimplementedUserServiceServer
	UserService service.UserService
}

func NewUserGrpcService(userService service.UserService) pbuser.UserServiceServer {
	return &UserGrpcServiceImplementation{
		UserService: userService,
	}
}

func (userGrpcService *UserGrpcServiceImplementation) Register(ctx context.Context, request *pbuser.RegisterRequest) (response *pbuser.RegisterResponse, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("cannot find metadata")
		helper.PrintLogToTerminal(err, "requestId")
		err = exception.GrpcErrorHandler(err)
		return
	}
	requestIdA, ok := md["requestid"]
	if !ok {
		err = errors.New("cannot finid requestid")
		helper.PrintLogToTerminal(err, "requestid")
		err = exception.GrpcErrorHandler(err)
		return
	}
	var registerUserRequest modelrequest.RegisterUserRequest
	registerUserRequest.Username = request.Username
	registerUserRequest.Email = request.Email
	registerUserRequest.Password = request.Password
	registerUserRequest.ConfirmPassword = request.Confirmpassword
	registerUserRequest.Utc = request.Utc
	registerUserResponse, err := userGrpcService.UserService.Register(ctx, requestIdA[0], registerUserRequest)
	if err != nil {
		err = exception.GrpcErrorHandler(err)
		return
	}

	registerResponse := pbuser.RegisterResponse{Username: registerUserResponse.Username, Email: registerUserResponse.Email, Utc: registerUserResponse.Utc}
	response = &registerResponse
	return
}

func (userGrpcService *UserGrpcServiceImplementation) Login(ctx context.Context, request *pbuser.LoginRequest) (response *pbuser.LoginResponse, err error) {
	var loginUserRequest modelrequest.LoginUserRequest
	loginUserRequest.Email = request.Email
	loginUserRequest.Password = request.Password
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("cannot find metadata")
		helper.PrintLogToTerminal(err, "requestId")
		err = exception.GrpcErrorHandler(err)
		return
	}
	requestIdA, ok := md["requestid"]
	if !ok {
		err = errors.New("cannot finid requestid")
		helper.PrintLogToTerminal(err, "requestid")
		err = exception.GrpcErrorHandler(err)
		return
	}
	sessionIdA, ok := md["sessionid"]
	if !ok {
		err = errors.New("cannot find sessionid")
		helper.PrintLogToTerminal(err, requestIdA[0])
		err = exception.GrpcErrorHandler(err)
		return
	}

	loginUserResponse, err := userGrpcService.UserService.Login(ctx, requestIdA[0], sessionIdA[0], loginUserRequest)
	if err != nil {
		err = exception.GrpcErrorHandler(err)
		return
	}
	loginResponse := pbuser.LoginResponse{Sessionid: loginUserResponse}
	response = &loginResponse
	return
}

func (userGrpcService *UserGrpcServiceImplementation) Logout(ctx context.Context, request *pbuser.LogoutRequest) (response *pbuser.LogoutResponse, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("cannot find metadata")
		helper.PrintLogToTerminal(err, "requestId")
		err = exception.GrpcErrorHandler(err)
		return
	}
	requestIdA, ok := md["requestid"]
	if !ok {
		err = errors.New("cannot finid requestid")
		helper.PrintLogToTerminal(err, "requestid")
		err = exception.GrpcErrorHandler(err)
		return
	}

	err = userGrpcService.UserService.Logout(ctx, requestIdA[0], request.GetSessionid())
	if err != nil {
		err = exception.GrpcErrorHandler(err)
		return
	}
	logoutResponse := pbuser.LogoutResponse{Msg: "successfully logout"}
	response = &logoutResponse
	return
}
