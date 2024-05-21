package routes

import (
	controller "gateway/controllers"

	"github.com/labstack/echo/v4"
)

func LandingPageRoute(e *echo.Echo, controller controller.LandingPageController) {
	e.GET("/", controller.LandingPage)
}
