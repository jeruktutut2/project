package setup

import (
	"context"
	"net/http"
	exception "project-user/exceptions"
	middleware "project-user/middlewares"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Echo(timeout uint8) (e *echo.Echo) {
	e = echo.New()
	e.Use(echomiddleware.Recover())
	e.HTTPErrorHandler = exception.CustomHTTPErrorHandler
	e.Use(middleware.SetGlobalRequestId)
	e.Use(middleware.SetGlobalTimeout(timeout))
	e.Use(middleware.SetGlobalRequestLog)
	return
}

func StartEcho(e *echo.Echo, port string) {
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	println(time.Now().String(), "echo: started")
}

func StopEcho(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cancel()
	println(time.Now().String(), "echo: shutdown properly")
}
