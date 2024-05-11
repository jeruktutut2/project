package controller

import (
	"net/http"
	exception "project-user/exceptions"
	helper "project-user/helpers"
	middleware "project-user/middlewares"
	modelrequest "project-user/models/requests"
	modelresponse "project-user/models/responses"
	service "project-user/services"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type UserControllerImplementation struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImplementation{
		UserService: userService,
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

	registerUserResponse, err := controller.UserService.Register(c.Request().Context(), requestId, registerUserRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	return modelresponse.ToResponse(c, http.StatusCreated, requestId, registerUserResponse, "")
}

func (controller *UserControllerImplementation) Login(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	loginUserRequest := modelrequest.LoginUserRequest{}
	if err := c.Bind(&loginUserRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}

	sessionId, err := controller.UserService.Login(c.Request().Context(), requestId, loginUserRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = sessionId
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully login", "")
}

func (controller *UserControllerImplementation) Logout(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	sessionId := c.Request().Context().Value(middleware.SessionIdKey).(string)
	err := controller.UserService.Logout(c.Request().Context(), requestId, sessionId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully logout", "")
}
