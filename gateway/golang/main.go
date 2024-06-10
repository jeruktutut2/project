package main

import (
	"context"
	"os"
	"os/signal"

	configuration "gateway/configurations"
	controller "gateway/controllers"
	route "gateway/routes"
	"gateway/setup"
)

func main() {
	config := configuration.NewConfiguration()

	// redisUtil := util.NewRedisConnection(config.RedisHost, config.RedisDatabase)
	// defer redisUtil.Close()

	e := setup.Echo(config.ProjectGatewayApplicationTimeout)
	setup.StartEcho(e, config.ProjectGatewayApplicationPort)
	defer setup.StopEcho(e)

	userClientConnection, userServiceClient := setup.NewUserClientConnection(config.ProjectGatewayUserApplicationHost)
	defer setup.CloseUserClientConnection(userClientConnection)
	userController := controller.NewUserController(userServiceClient)
	landingPageController := controller.NewLandingPageController()
	route.UserRoute(e, userController)
	route.LandingPageRoute(e, landingPageController)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
