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
	config := configuration.NewConfiguration()

	mysqlUtil := util.NewMysqlConnection(config.MysqlUsername, config.MysqlPassword, config.MysqlHost, config.MysqlDatabase, config.MysqlMaxOpenConnection, config.MysqlMaxIdleConnection, config.MysqlConnectionMaxLifetime, config.MysqlConnectionMaxIdletime)
	defer mysqlUtil.Close()

	redisUtil := util.NewRedisConnection(config.RedisHost, config.RedisDatabase)
	defer redisUtil.Close()

	validate := setup.Validate()

	repositorySetup := setup.NewRepositorySetup()
	serviceSetup := setup.NewServiceSetup(mysqlUtil, redisUtil, validate, repositorySetup)
	grpcSetup := setup.NewGrpcSetup(serviceSetup, config.ApplicationPort)
	defer setup.StopGrpc(grpcSetup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
