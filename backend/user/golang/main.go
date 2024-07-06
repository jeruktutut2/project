package main

import (
	"context"
	"os"
	"os/signal"
	configuration "project-user/configurations"
	"project-user/helpers"
	repository "project-user/repositories"
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

	bcryptHelper := helpers.NewBcryptHelper()
	timeHelper := helpers.NewTimeHelper()
	redisRepository := repository.NewRedisRepository()
	uuidHelper := helpers.NewUuidHelper()
	repositorySetup := setup.NewRepositorySetup()
	serviceSetup := setup.NewServiceSetup(mysqlUtil, redisUtil, validate, repositorySetup, bcryptHelper, timeHelper, redisRepository, uuidHelper)
	grpcSetup := setup.NewUserGrpcSetup(serviceSetup, config.ProjectUserApplicationHost)
	defer setup.StopUserGrpc(grpcSetup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
