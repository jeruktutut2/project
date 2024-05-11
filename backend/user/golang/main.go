package main

import (
	"context"
	"os"
	"os/signal"
	configuration "project-user/configurations"
	"project-user/setup"
	util "project-user/utils"
)

func main() {
	// fmt.Println(time.Now().String())
	config := configuration.NewConfiguration()

	mysqlUtil := util.NewMysqlConnection(config.MysqlUsername, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase, config.MysqlMaxOpenConnection, config.MysqlMaxIdleConnection, config.MysqlConnectionMaxLifetime, config.MysqlConnectionMaxIdletime)
	defer mysqlUtil.Close()

	redisUtil := util.NewRedisConnection(config.RedisHost, config.RedisPort, config.RedisDatabase)
	defer redisUtil.Close()

	e := setup.Echo(config.ApplicationTimeout)
	setup.StartEcho(e, config.ApplicationPort)

	validate := setup.Validate()

	repositorySetup := setup.NewRepositorySetup()
	serviceSetup := setup.NewServiceSetup(mysqlUtil, redisUtil, validate, repositorySetup)
	controllerSetup := setup.NewControllerSetup(serviceSetup)
	setup.RouteSetup(e, redisUtil, controllerSetup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	setup.StopEcho(e)
}
