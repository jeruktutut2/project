package setup

import (
	"project-user/helpers"
	repository "project-user/repositories"
	service "project-user/services"
	util "project-user/utils"

	"github.com/go-playground/validator/v10"
)

type ServiceSetup struct {
	UserService service.UserService
}

func NewServiceSetup(mysqlUtil util.MysqlUtil, redisUtil util.RedisUtil, validate *validator.Validate, repositorySetup *RepositorySetup, bcryptHelper helpers.BcryptHelper, timeHelper helpers.TimeHelper, redisRepository repository.RedisRepository, uuidHelper helpers.UuidHelper) *ServiceSetup {
	return &ServiceSetup{
		UserService: service.NewUserService(mysqlUtil, redisUtil, validate, repositorySetup.UserRepository, repositorySetup.UserPermissionRepository, bcryptHelper, timeHelper, redisRepository, uuidHelper),
	}
}
