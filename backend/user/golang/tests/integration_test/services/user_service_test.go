package services_test

import (
	"context"
	"project-user/helpers"
	modelrequests "project-user/models/requests"
	modelresponses "project-user/models/responses"
	repositories "project-user/repositories"
	services "project-user/services"
	"project-user/tests/initialize"
	"project-user/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx                      context.Context
	requestId                string
	sessionId                string
	registerUserRequest      modelrequests.RegisterUserRequest
	loginUserRequest         modelrequests.LoginUserRequest
	mysqlUtil                utils.MysqlUtil
	redisUtil                utils.RedisUtil
	validate                 *validator.Validate
	userRepository           repositories.UserRepository
	userPermissionRepository repositories.UserPermissionRepository
	bcryptHelper             helpers.BcryptHelper
	timeHelper               helpers.TimeHelper
	redisRepository          repositories.RedisRepository
	uuidHelper               helpers.UuidHelper
	userService              services.UserService
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (sut *UserServiceTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.mysqlUtil = utils.NewMysqlConnection("root", "12345", "localhost:3309", "user", 10, 10, 10, 10)
	sut.redisUtil = utils.NewRedisConnection("localhost", "6380", 0)
	sut.validate = validator.New()
	helpers.UsernameValidator(sut.validate)
	helpers.PasswordValidator(sut.validate)
	helpers.TelephoneValidator(sut.validate)
	sut.userRepository = repositories.NewUserRepository()
	sut.userPermissionRepository = repositories.NewUserPermissinoRepository()
	sut.bcryptHelper = helpers.NewBcryptHelper()
	sut.timeHelper = helpers.NewTimeHelper()
	sut.redisRepository = repositories.NewRedisRepository()
	sut.uuidHelper = helpers.NewUuidHelper()
	sut.userService = services.NewUserService(sut.mysqlUtil, sut.redisUtil, sut.validate, sut.userRepository, sut.userPermissionRepository, sut.bcryptHelper, sut.timeHelper, sut.redisRepository, sut.uuidHelper)
}

func (sut *UserServiceTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.ctx = context.Background()
	sut.requestId = "requestId"
	sut.sessionId = "sessionId"
	sut.registerUserRequest = modelrequests.RegisterUserRequest{
		Username:        "username",
		Email:           "email@email.com",
		Password:        "password@A1",
		Confirmpassword: "password@A1",
		Utc:             "+0800",
	}
	sut.loginUserRequest = modelrequests.LoginUserRequest{
		Email:    "email@email.com",
		Password: "password@A1",
	}
}

func (sut *UserServiceTestSuite) BeforeTest(suiteName, testName string) {
	sut.T().Log("BeforeTest: " + suiteName + " " + testName)
}

func (sut *UserServiceTestSuite) TestRegisterRegisterUserRequestValidationError() {
	sut.T().Log("TestRegisterRegisterUserRequestValidationError")
	registerUserRequest := modelrequests.RegisterUserRequest{}
	response, err := sut.userService.Register(sut.ctx, sut.requestId, registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"username","message":"is required"},{"field":"email","message":"is required"},{"field":"password","message":"is required"},{"field":"confirmpassword","message":"is required"},{"field":"utc","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterPasswordAndConfirmpasswordIsDifferent() {
	sut.T().Log("TestRegisterPasswordAndConfirmpasswordIsDifferent")
	sut.registerUserRequest.Confirmpassword = "password@A1-"
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"password and confirm password is different"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameInternalServerError() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameBadRequestUsernameAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameBadRequestUsernameAlreadyExists")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"username already exists"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByCountByEmailBadRequestUsernameAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameBadRequestUsernameAlreadyExists")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.registerUserRequest.Username = "username1"
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"email already exists"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterSuccess() {
	sut.T().Log("TestRegisterSuccess")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.registerUserRequest.Username = "username1"
	sut.registerUserRequest.Email = "email1@email.com"
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	registerUserResponse := modelresponses.RegisterUserResponse{}
	registerUserResponse.Username = "username1"
	registerUserResponse.Email = "email1@email.com"
	registerUserResponse.Utc = "+0800"
	sut.Equal(response, registerUserResponse)
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithSessionIdInternalServerError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithSessionIdInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.loginUserRequest = modelrequests.LoginUserRequest{}
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailInternalServerError() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.loginUserRequest.Password = "password@A1-"
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("TestLoginUserPermissionRepositoryFindByUserIdInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLoginSuccess() {
	sut.T().Log("TestLoginSuccess")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.NotEqual(response, "")
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLogoutRowsAffectedNotOneInternalServerError() {
	sut.T().Log("TestLogoutRowsAffectedNotOneInternalServerError")
	err := sut.userService.Logout(sut.ctx, sut.requestId, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLogoutSuccess() {
	sut.T().Log("TestLogoutSuccess")
	initialize.DelDataRedis(sut.redisUtil.GetClient(), sut.ctx, sut.sessionId)
	initialize.SetDataRedis(sut.redisUtil.GetClient(), sut.ctx, sut.sessionId, "value", time.Duration(0))
	err := sut.userService.Logout(sut.ctx, sut.requestId, sut.sessionId)
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) AfterTest(suiteName, testName string) {
	sut.T().Log("AfterTest: " + suiteName + " " + testName)
}

func (sut *UserServiceTestSuite) TearDownTest() {
	sut.T().Log("TearDownTest")
}

func (sut *UserServiceTestSuite) TearDownSuite() {
	sut.T().Log("TearDownSuite")
}
