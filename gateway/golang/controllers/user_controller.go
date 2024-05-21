package controller

import (
	exception "gateway/exceptions"
	helper "gateway/helpers"
	middleware "gateway/middlewares"
	modelrequest "gateway/models/requests"
	modelresponse "gateway/models/responses"
	pbuser "gateway/protofiles/pb/api/v1/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

type UserController interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type UserControllerImplementation struct {
	UserServiceClient pbuser.UserServiceClient
}

func NewUserController(userServiceclient pbuser.UserServiceClient) UserController {
	return &UserControllerImplementation{
		UserServiceClient: userServiceclient,
	}
}

func (controller *UserControllerImplementation) Register(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	registerUserRequest := modelrequest.RegisterUserRequest{}
	if err := c.Bind(&registerUserRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	registerRequest := &pbuser.RegisterRequest{
		Username:        registerUserRequest.Username,
		Email:           registerUserRequest.Email,
		Password:        registerUserRequest.Password,
		Confirmpassword: registerUserRequest.ConfirmPassword,
		Utc:             registerUserRequest.Utc,
	}

	md := metadata.Pairs(
		"requestid", requestId,
	)
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
	registerResponse, err := controller.UserServiceClient.Register(ctx, registerRequest)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		return exception.GrpcErrorHandler(c, requestId, err)
	}

	response := modelresponse.ToRegisterUserResponse(registerResponse)

	return modelresponse.ToResponse(c, http.StatusCreated, requestId, response, "")
}

func (controller *UserControllerImplementation) Login(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	sessionId := c.Request().Context().Value(middleware.SessionIdKey).(string)
	var loginUserRequest modelrequest.LoginUserRequest
	if err := c.Bind(&loginUserRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	loginRequest := &pbuser.LoginRequest{Email: loginUserRequest.Email, Password: loginUserRequest.Password}
	md := metadata.Pairs(
		// why lowercase, because it will automatically lowercase send to server
		"requestid", requestId,
		"sessionid", sessionId,
	)
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
	loginResponse, err := controller.UserServiceClient.Login(ctx, loginRequest)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		return exception.GrpcErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = loginResponse.GetSessionid()
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully login", "")
}

func (controller *UserControllerImplementation) Logout(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully logout", "")
}
