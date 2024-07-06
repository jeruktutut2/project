package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	exception "project-user/exceptions"
	"project-user/helpers"
	helper "project-user/helpers"
	modelentities "project-user/models/entities"
	modelrequests "project-user/models/requests"
	modelresponse "project-user/models/responses"
	repository "project-user/repositories"
	utils "project-user/utils"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, requestId string, registerUserRequest modelrequests.RegisterUserRequest) (registerUserResponse modelresponse.RegisterUserResponse, err error)
	Login(ctx context.Context, requestId string, sessionIdUser string, loginUserRequest modelrequests.LoginUserRequest) (sessionId string, err error)
	Logout(ctx context.Context, requestId string, sessionId string) (err error)
}

type UserServiceImplementation struct {
	MysqlUtil                utils.MysqlUtil
	RedisUtil                utils.RedisUtil
	Validate                 *validator.Validate
	UserRepository           repository.UserRepository
	UserPermissionRepository repository.UserPermissionRepository
	BcryptHelper             helpers.BcryptHelper
	TimeHelper               helpers.TimeHelper
	RedisRepository          repository.RedisRepository
	UuidHelper               helper.UuidHelper
}

func NewUserService(mysqlUtil utils.MysqlUtil, redisUtil utils.RedisUtil, validate *validator.Validate, userRepository repository.UserRepository, userPermissionRepository repository.UserPermissionRepository, bcryptHelper helpers.BcryptHelper, timeHelper helpers.TimeHelper, redisRepository repository.RedisRepository, uuidHelper helper.UuidHelper) UserService {
	return &UserServiceImplementation{
		MysqlUtil:                mysqlUtil,
		RedisUtil:                redisUtil,
		Validate:                 validate,
		UserRepository:           userRepository,
		UserPermissionRepository: userPermissionRepository,
		BcryptHelper:             bcryptHelper,
		TimeHelper:               timeHelper,
		RedisRepository:          redisRepository,
		UuidHelper:               uuidHelper,
	}
}

func (service *UserServiceImplementation) Register(ctx context.Context, requestId string, registerUserRequest modelrequests.RegisterUserRequest) (registerUserResponse modelresponse.RegisterUserResponse, err error) {
	err = service.Validate.Struct(registerUserRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, registerUserRequest)
		if validationResult != nil {
			err = exception.ToResponseExceptionRequestValidation(requestId, validationResult)
			return
		}
	}
	if registerUserRequest.Password != registerUserRequest.Confirmpassword {
		err = errors.New("password and confirm password is different")
		err = exception.ToResponseException(err, requestId, http.StatusBadRequest, "password and confirm password is different")
		return
	}

	numberOfUser, err := service.UserRepository.CountByUsername(service.MysqlUtil.GetDb(), ctx, registerUserRequest.Username)
	if err != nil && err != sql.ErrNoRows {
		err = exception.CheckError(err, requestId)
		return
	}
	if numberOfUser > 0 {
		err = errors.New("username already exists")
		err = exception.ToResponseException(err, requestId, http.StatusBadRequest, "username already exists")
		return
	}

	numberOfUser, err = service.UserRepository.CountByEmail(service.MysqlUtil.GetDb(), ctx, registerUserRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		err = exception.CheckError(err, requestId)
		return
	}
	if numberOfUser > 0 {
		err = errors.New("email already exists")
		err = exception.ToResponseException(err, requestId, http.StatusBadRequest, "email already exists")
		return
	}

	hashedPassword, err := service.BcryptHelper.GenerateFromPassword([]byte(registerUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		err = exception.CheckError(err, requestId)
		return
	}

	var user modelentities.User
	user.Username = sql.NullString{Valid: true, String: registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: service.TimeHelper.NowUnixMili()}
	rowsAffected, err := service.UserRepository.Create(service.MysqlUtil.GetDb(), ctx, user)
	if err != nil {
		err = exception.CheckError(err, requestId)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not 1 when create user")
		err = exception.CheckError(err, requestId)
		return
	}

	registerUserResponse = modelresponse.ToRegisterUserResponse(user)
	return
}

func (service *UserServiceImplementation) Login(ctx context.Context, requestId string, sessionIdUser string, loginUserRequest modelrequests.LoginUserRequest) (sessionId string, err error) {
	if sessionIdUser != "" {
		var rowsAffected int64
		rowsAffected, err = service.RedisRepository.Del(service.RedisUtil.GetClient(), ctx, sessionIdUser)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
		} else if err == nil && rowsAffected != 1 {
			err = errors.New("rows affected not 1")
			helper.PrintLogToTerminal(err, requestId)
		}
	}
	err = service.Validate.Struct(loginUserRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, loginUserRequest)
		if validationResult != nil {
			err = exception.ToResponseExceptionRequestValidation(requestId, validationResult)
			return
		}
	}

	user, err := service.UserRepository.FindByEmail(service.MysqlUtil.GetDb(), ctx, loginUserRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		err = exception.CheckError(err, requestId)
		return
	} else if err != nil && err == sql.ErrNoRows {
		err = errors.New("wrong email or password")
		err = exception.ToResponseException(err, requestId, http.StatusBadRequest, "wrong email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginUserRequest.Password))
	if err != nil {
		err = errors.New("wrong email or password")
		err = exception.ToResponseException(err, requestId, http.StatusBadRequest, "wrong email or password")
		return
	}

	userPermissions, err := service.UserPermissionRepository.FindByUserId(service.MysqlUtil.GetDb(), ctx, user.Id.Int32)
	if err != nil {
		err = exception.CheckError(err, requestId)
		return
	}
	var idPermissions []int32
	for _, userPermission := range userPermissions {
		idPermissions = append(idPermissions, userPermission.PermissionId.Int32)
	}

	sessionId = service.UuidHelper.String()
	sessionValue := make(map[string]interface{})
	sessionValue["id"] = user.Id.Int32
	sessionValue["username"] = user.Username.String
	sessionValue["email"] = user.Email.String
	sessionValue["idPermissions"] = idPermissions
	sessionByte, err := json.Marshal(sessionValue)
	if err != nil {
		err = exception.CheckError(err, requestId)
		return
	}
	session := string(sessionByte)

	_, err = service.RedisRepository.Set(service.RedisUtil.GetClient(), ctx, sessionId, session, 0)
	if err != nil && err != redis.Nil {
		err = exception.CheckError(err, requestId)
		return
	}
	return
}

func (service *UserServiceImplementation) Logout(ctx context.Context, requestId string, sessionId string) (err error) {
	rowsAffected, err := service.RedisRepository.Del(service.RedisUtil.GetClient(), ctx, sessionId)
	if err != nil && err != redis.Nil {
		err = exception.CheckError(err, requestId)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not 1 when delete data to redis")
		err = exception.CheckError(err, requestId)
		return
	}
	return
}
