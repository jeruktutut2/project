package services_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"project-user/helpers"
	modelentities "project-user/models/entities"
	modelrequests "project-user/models/requests"
	modelresponses "project-user/models/responses"
	service "project-user/services"
	mockhelpers "project-user/tests/unit_tests/mocks/helpers"
	mockrepositories "project-user/tests/unit_tests/mocks/repositories"
	mockutils "project-user/tests/unit_tests/mocks/utils"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx                          context.Context
	requestId                    string
	sessionId                    string
	db                           *sql.DB
	options                      *sql.TxOptions
	tx                           *sql.Tx
	client                       *redis.Client
	errTimeout                   error
	errInternalServer            error
	registerUserRequest          modelrequests.RegisterUserRequest
	loginUserRequest             modelrequests.LoginUserRequest
	user                         modelentities.User
	mysqlUtilMock                *mockutils.MysqlUtilMock
	redisUtilMock                *mockutils.RedisUtilMock
	validate                     *validator.Validate
	userRepositoryMock           *mockrepositories.UserRepositoryMock
	userPermissionRepositoryMock *mockrepositories.UserPermissionRepositoryMock
	bcryptHelperMock             *mockhelpers.BcryptHelperMock
	timeHelperMock               *mockhelpers.TimeHelperMock
	redisRepositoryMock          *mockrepositories.RedisRepositoryMock
	uuidHelperMock               *mockhelpers.UuidHelperMock
	userService                  service.UserService
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (sut *UserServiceTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.ctx = context.Background()
	sut.requestId = "requestId"
	sut.sessionId = "sessionId"
	sut.db = &sql.DB{}
	sut.options = &sql.TxOptions{}
	sut.tx = &sql.Tx{}
	sut.client = &redis.Client{}
	sut.errTimeout = context.Canceled
	sut.errInternalServer = errors.New("internal server error")
}

func (sut *UserServiceTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.registerUserRequest = modelrequests.RegisterUserRequest{
		Username:        "username17",
		Email:           "email17@email.com",
		Password:        "password@A1",
		Confirmpassword: "password@A1",
		Utc:             "+0800",
	}
	sut.loginUserRequest = modelrequests.LoginUserRequest{
		Email:    "email17@email.com",
		Password: "password@A1",
	}
	sut.user = modelentities.User{
		Id:        sql.NullInt32{Valid: true, Int32: 1},
		Username:  sql.NullString{Valid: true, String: "username17"},
		Email:     sql.NullString{Valid: true, String: "email17@email.com"},
		Password:  sql.NullString{Valid: true, String: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG"},
		Utc:       sql.NullString{Valid: true, String: "+0800"},
		CreatedAt: sql.NullInt64{Valid: true, Int64: 1719496855216},
	}
	sut.mysqlUtilMock = new(mockutils.MysqlUtilMock)
	sut.redisUtilMock = new(mockutils.RedisUtilMock)
	sut.validate = validator.New()
	helpers.UsernameValidator(sut.validate)
	helpers.PasswordValidator(sut.validate)
	helpers.TelephoneValidator(sut.validate)
	sut.userRepositoryMock = new(mockrepositories.UserRepositoryMock)
	sut.userPermissionRepositoryMock = new(mockrepositories.UserPermissionRepositoryMock)
	sut.bcryptHelperMock = new(mockhelpers.BcryptHelperMock)
	sut.timeHelperMock = new(mockhelpers.TimeHelperMock)
	sut.redisRepositoryMock = new(mockrepositories.RedisRepositoryMock)
	sut.uuidHelperMock = new(mockhelpers.UuidHelperMock)
	sut.userService = service.NewUserService(sut.mysqlUtilMock, sut.redisUtilMock, sut.validate, sut.userRepositoryMock, sut.userPermissionRepositoryMock, sut.bcryptHelperMock, sut.timeHelperMock, sut.redisRepositoryMock, sut.uuidHelperMock)
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

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameTimeoutError() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameTimeoutError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, sut.errTimeout)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameInternalServerError() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameInternalServerError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, sut.errInternalServer)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameUsernameAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameUsernameAlreadyExists")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(1, nil)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"username already exists"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByEmailTimeoutError() {
	sut.T().Log("TestRegisterUserRepositoryCountByEmailTimeoutError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, sut.errTimeout)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByEmailInternalServerError() {
	sut.T().Log("TestRegisterUserRepositoryCountByEmailInternalServerError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, sut.errInternalServer)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByEmailEmailAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByEmailEmailAlreadyExists")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(1, nil)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"email already exists"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterBcryptGenerateFromPasswordTimeoutError() {
	sut.T().Log("TestRegisterBcryptGenerateFromPasswordTimeoutError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return([]uint8{}, sut.errTimeout)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterBcryptGenerateFromPasswordInternalServerError() {
	sut.T().Log("TestRegisterBcryptGenerateFromPasswordInternalServerError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return([]uint8{}, sut.errInternalServer)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCreateTimeoutError() {
	sut.T().Log("TestRegisterUserRepositoryCreateTimeoutError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	hashedPassword := []uint8{36, 50, 97, 36, 49, 48, 36, 121, 50, 52, 100, 109, 102, 86, 117, 108, 116, 68, 115, 118, 56, 55, 97, 77, 52, 105, 52, 89, 101, 76, 73, 74, 70, 118, 57, 104, 85, 54, 104, 119, 112, 97, 112, 65, 82, 102, 104, 84, 103, 121, 48, 116, 111, 75, 72, 78, 108, 114, 76, 83}
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)
	var nowUnixMili int64 = 1719496855216
	sut.timeHelperMock.Mock.On("NowUnixMili").Return(nowUnixMili)
	var user modelentities.User
	user.Username = sql.NullString{Valid: true, String: sut.registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: sut.registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: sut.registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: nowUnixMili}
	sut.userRepositoryMock.Mock.On("Create", sut.db, sut.ctx, user).Return(int64(0), sut.errTimeout)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCreateInternalServerError() {
	sut.T().Log("TestRegisterUserRepositoryCreateInternalServerError")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	hashedPassword := []uint8{36, 50, 97, 36, 49, 48, 36, 121, 50, 52, 100, 109, 102, 86, 117, 108, 116, 68, 115, 118, 56, 55, 97, 77, 52, 105, 52, 89, 101, 76, 73, 74, 70, 118, 57, 104, 85, 54, 104, 119, 112, 97, 112, 65, 82, 102, 104, 84, 103, 121, 48, 116, 111, 75, 72, 78, 108, 114, 76, 83}
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)
	var nowUnixMili int64 = 1719496855216
	sut.timeHelperMock.Mock.On("NowUnixMili").Return(nowUnixMili)
	var user modelentities.User
	user.Username = sql.NullString{Valid: true, String: sut.registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: sut.registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: sut.registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: nowUnixMili}
	sut.userRepositoryMock.Mock.On("Create", sut.db, sut.ctx, user).Return(int64(0), sut.errInternalServer)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCreateRowsAffectedNotOne() {
	sut.T().Log("TestRegisterUserRepositoryCreateRowsAffectedNotOne")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	hashedPassword := []uint8{36, 50, 97, 36, 49, 48, 36, 121, 50, 52, 100, 109, 102, 86, 117, 108, 116, 68, 115, 118, 56, 55, 97, 77, 52, 105, 52, 89, 101, 76, 73, 74, 70, 118, 57, 104, 85, 54, 104, 119, 112, 97, 112, 65, 82, 102, 104, 84, 103, 121, 48, 116, 111, 75, 72, 78, 108, 114, 76, 83}
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)
	var nowUnixMili int64 = 1719496855216
	sut.timeHelperMock.Mock.On("NowUnixMili").Return(nowUnixMili)
	var user modelentities.User
	user.Username = sql.NullString{Valid: true, String: sut.registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: sut.registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: sut.registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: nowUnixMili}
	sut.userRepositoryMock.Mock.On("Create", sut.db, sut.ctx, user).Return(int64(0), nil)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	sut.Equal(response, modelresponses.RegisterUserResponse{})
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterSuccess() {
	sut.T().Log("TestRegisterSuccess")
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("CountByUsername", sut.db, sut.ctx, sut.registerUserRequest.Username).Return(0, nil)
	sut.userRepositoryMock.Mock.On("CountByEmail", sut.db, sut.ctx, sut.registerUserRequest.Email).Return(0, nil)
	hashedPassword := []uint8{36, 50, 97, 36, 49, 48, 36, 121, 50, 52, 100, 109, 102, 86, 117, 108, 116, 68, 115, 118, 56, 55, 97, 77, 52, 105, 52, 89, 101, 76, 73, 74, 70, 118, 57, 104, 85, 54, 104, 119, 112, 97, 112, 65, 82, 102, 104, 84, 103, 121, 48, 116, 111, 75, 72, 78, 108, 114, 76, 83}
	sut.bcryptHelperMock.Mock.On("GenerateFromPassword", []byte(sut.registerUserRequest.Password), bcrypt.DefaultCost).Return(hashedPassword, nil)
	var nowUnixMili int64 = 1719496855216
	sut.timeHelperMock.Mock.On("NowUnixMili").Return(nowUnixMili)
	var user modelentities.User
	user.Username = sql.NullString{Valid: true, String: sut.registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: sut.registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: sut.registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: nowUnixMili}
	sut.userRepositoryMock.Mock.On("Create", sut.db, sut.ctx, user).Return(int64(1), nil)
	response, err := sut.userService.Register(sut.ctx, sut.requestId, sut.registerUserRequest)
	registerUserResponse := modelresponses.RegisterUserResponse{}
	registerUserResponse.Username = "username17"
	registerUserResponse.Email = "email17@email.com"
	registerUserResponse.Utc = "+0800"
	sut.Equal(response, registerUserResponse)
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithSessionIdTimeoutError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithSessionIdTimeoutError")
	sut.loginUserRequest = modelrequests.LoginUserRequest{}
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errTimeout)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithSessionIdInternalServerError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithSessionIdInternalServerError")
	sut.loginUserRequest = modelrequests.LoginUserRequest{}
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, sut.sessionId, sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithoutSessionIdTimeoutError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithoutSessionIdTimeoutError")
	sut.loginUserRequest = modelrequests.LoginUserRequest{}
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errTimeout)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithoutSessionIdInternalServerError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithoutSessionIdInternalServerError")
	sut.loginUserRequest = modelrequests.LoginUserRequest{}
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailTimeoutError() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailTimeoutError")

	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(modelentities.User{}, sut.errTimeout)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailInternalServerError() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailInternalServerError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(modelentities.User{}, sut.errInternalServer)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailBadRequestWrongEmailPassword() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailBadRequestWrongEmailPassword")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(modelentities.User{}, sql.ErrNoRows)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailPassword() {
	sut.T().Log("TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailPassword")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.user.Password.String = ""
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserPermissionRepositoryFindByUserIdTimeoutError() {
	sut.T().Log("TestLoginUserPermissionRepositoryFindByUserIdTimeoutError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.db, sut.ctx, sut.user.Id.Int32).Return([]modelentities.UserPermission{}, sut.errTimeout)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("TestLoginUserPermissionRepositoryFindByUserIdInternalServerError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.db, sut.ctx, sut.user.Id.Int32).Return([]modelentities.UserPermission{}, sut.errInternalServer)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, "")
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositorySetTimeoutError() {
	sut.T().Log("TestLoginRedisRepositorySetTimeoutError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.db, sut.ctx, sut.user.Id.Int32).Return([]modelentities.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email17@email.com","id":1,"idPermissions":null,"username":"username17"}`
	sut.redisRepositoryMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", sut.errTimeout)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositorySetInternalServerError() {
	sut.T().Log("TestLoginRedisRepositorySetInternalServerError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.db, sut.ctx, sut.user.Id.Int32).Return([]modelentities.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email17@email.com","id":1,"idPermissions":null,"username":"username17"}`
	sut.redisRepositoryMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", sut.errInternalServer)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	fmt.Println("response:", response, err)
	sut.Equal(response, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLoginSuccess() {
	sut.T().Log("TestLoginSuccess")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	sut.mysqlUtilMock.Mock.On("GetDb").Return(sut.db)
	sut.userRepositoryMock.Mock.On("FindByEmail", sut.db, sut.ctx, sut.loginUserRequest.Email).Return(sut.user, nil)
	sut.userPermissionRepositoryMock.Mock.On("FindByUserId", sut.db, sut.ctx, sut.user.Id.Int32).Return([]modelentities.UserPermission{}, nil)
	sut.uuidHelperMock.Mock.On("String").Return(sut.sessionId)
	session := `{"email":"email17@email.com","id":1,"idPermissions":null,"username":"username17"}`
	sut.redisRepositoryMock.Mock.On("Set", sut.client, sut.ctx, sut.sessionId, session, time.Duration(0)).Return("", nil)
	response, err := sut.userService.Login(sut.ctx, sut.requestId, "", sut.loginUserRequest)
	sut.Equal(response, sut.sessionId)
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLogoutRedisRepositoryDelTimeoutError() {
	sut.T().Log("TestLogoutRedisRepositoryDelTimeoutError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errTimeout)
	err := sut.userService.Logout(sut.ctx, sut.requestId, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"time out or user cancel the request"}]`)
}

func (sut *UserServiceTestSuite) TestLogoutRedisRepositoryDelInternalServerError() {
	sut.T().Log("TestLogoutRedisRepositoryDelInternalServerError")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), sut.errInternalServer)
	err := sut.userService.Logout(sut.ctx, sut.requestId, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLogoutRedisRepositoryDelRowsAffectedNotOne() {
	sut.T().Log("TestLogoutRedisRepositoryDelRowsAffectedNotOne")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(0), nil)
	err := sut.userService.Logout(sut.ctx, sut.requestId, sut.sessionId)
	sut.Equal(err.Error(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLogoutSuccess() {
	sut.T().Log("TestLogoutSuccess")
	sut.redisUtilMock.Mock.On("GetClient").Return(sut.client)
	sut.redisRepositoryMock.Mock.On("Del", sut.client, sut.ctx, sut.sessionId).Return(int64(1), nil)
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
