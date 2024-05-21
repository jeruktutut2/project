package setup

import (
	service "project-user/services"
	util "project-user/utils"

	"github.com/go-playground/validator/v10"
)

type ServiceSetup struct {
	UserService service.UserService
}

func NewServiceSetup(mysqlUtil util.MysqlUtil, redisUtil util.RedisUtil, validate *validator.Validate, repositorySetup *RepositorySetup) *ServiceSetup {
	return &ServiceSetup{
		UserService: service.NewUserService(mysqlUtil, redisUtil, validate, repositorySetup.UserRepository, repositorySetup.UserPermissionRepository),
	}
}
