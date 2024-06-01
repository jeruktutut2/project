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

	mysqlUtil := util.NewMysqlConnection(config.ProjectUserMysqlUsername, config.ProjectUserMysqlPassword, config.ProjectUserMysqlHost, config.ProjectUserMysqlDatabase, config.ProjectUserMysqlMaxOpenConnection, config.ProjectUserMysqlMaxIdleConnection, config.ProjectUserMysqlConnectionMaxLifetime, config.ProjectUserMysqlConnectionMaxIdletime)
	defer mysqlUtil.Close()

	redisUtil := util.NewRedisConnection(config.ProjectUserRedisHost, config.ProjectUserRedisPort, config.ProjectUserRedisDatabase)
	defer redisUtil.Close()

	validate := setup.Validate()

	repositorySetup := setup.NewRepositorySetup()
	serviceSetup := setup.NewServiceSetup(mysqlUtil, redisUtil, validate, repositorySetup)
	grpcSetup := setup.NewGrpcSetup(serviceSetup, config.ProjectUserApplicationPort)
	defer setup.StopGrpc(grpcSetup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
