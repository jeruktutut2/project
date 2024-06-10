package setup

import (
	"context"
	exception "gateway/exceptions"
	middleware "gateway/middlewares"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Echo(timeout int) (e *echo.Echo) {
	e = echo.New()
	e.Use(middleware.SetGlobalRequestId)
	e.Use(middleware.SetGlobalTimeout(timeout))
	e.Use(middleware.SetGlobalRequestLog)
	e.Use(echomiddleware.Recover())
	e.HTTPErrorHandler = exception.CustomHTTPErrorHandler
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
