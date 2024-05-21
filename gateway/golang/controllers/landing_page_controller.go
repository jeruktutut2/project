package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type LandingPageController interface {
	LandingPage(c echo.Context) error
}

type LandingPageControllerImplementation struct {
}

func NewLandingPageController() LandingPageController {
	return &LandingPageControllerImplementation{}
}

func (controller *LandingPageControllerImplementation) LandingPage(c echo.Context) error {
	s := struct {
		Name string `json:"name"`
		Date string `json:"date"`
	}{
		Name: "laning page",
		Date: time.Now().String(),
	}
	return c.JSON(http.StatusOK, s)
}
