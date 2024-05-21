package routes

import (
	controller "gateway/controllers"
	middleware "gateway/middlewares"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo, controller controller.UserController) {
	e.POST("/api/v1/user/register", controller.Register)
	e.POST("/api/v1/user/login", controller.Login, middleware.GetSessionId)
	e.POST("/api/v1/user/logout", controller.Logout)
}
