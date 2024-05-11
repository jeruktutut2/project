package route

import (
	controller "project-user/controllers"
	middleware "project-user/middlewares"
	util "project-user/utils"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo, redisUtil util.RedisUtil, controller controller.UserController) {
	e.POST("/api/v1/user/register", controller.Register)
	e.POST("/api/v1/user/login", controller.Login)
	e.POST("/api/v1/user/logout", controller.Logout, middleware.Authenticate(redisUtil))
}
