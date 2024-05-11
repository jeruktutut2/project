package setup

import (
	route "project-user/routes"
	util "project-user/utils"

	"github.com/labstack/echo/v4"
)

func RouteSetup(e *echo.Echo, redisUtil util.RedisUtil, controllerSetup ControllerSetup) {
	route.UserRoute(e, redisUtil, controllerSetup.UserController)
}
