package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	exception "project-user/exceptions"
	helper "project-user/helpers"
	modelentity "project-user/models/entities"
	modelrequest "project-user/models/requests"
	modelresponse "project-user/models/responses"
	repository "project-user/repositories"
	util "project-user/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, requestId string, registerUserRequest modelrequest.RegisterUserRequest) (registerUserResponse modelresponse.RegisterUserResponse, err error)
	Login(ctx context.Context, requestId string, sessionIdUser string, loginUserRequest modelrequest.LoginUserRequest) (sessionId string, err error)
	Logout(ctx context.Context, requestId string, sessionId string) (err error)
}

type UserServiceImplementation struct {
	MysqlUtil                util.MysqlUtil
	RedisUtil                util.RedisUtil
	Validate                 *validator.Validate
	UserRepository           repository.UserRepository
	UserPermissionRepository repository.UserPermissionRepository
}

func NewUserService(mysqlUtil util.MysqlUtil, redisUtil util.RedisUtil, validate *validator.Validate, userRepository repository.UserRepository, userPermissionRepository repository.UserPermissionRepository) UserService {
	return &UserServiceImplementation{
		MysqlUtil:                mysqlUtil,
		RedisUtil:                redisUtil,
		Validate:                 validate,
		UserRepository:           userRepository,
		UserPermissionRepository: userPermissionRepository,
	}
}

func (service *UserServiceImplementation) Register(ctx context.Context, requestId string, registerUserRequest modelrequest.RegisterUserRequest) (registerUserResponse modelresponse.RegisterUserResponse, err error) {
	err = service.Validate.Struct(registerUserRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, registerUserRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}
	if registerUserRequest.Password != registerUserRequest.ConfirmPassword {
		var results []helper.Result
		var result helper.Result
		result.Field = "password and confirmPassword"
		result.Message = "password and confirm password is different"
		results = append(results, result)
		var resultsByte []byte
		resultsByte, err = json.Marshal(results)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}
		err = exception.NewValidationException(string(resultsByte))
		return
	}

	numberOfUser, err := service.UserRepository.CountByUsername(service.MysqlUtil.GetDb(), ctx, registerUserRequest.Username)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if numberOfUser > 0 {
		err = exception.NewBadRequestException("username already exists")
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	numberOfUser, err = service.UserRepository.CountByEmail(service.MysqlUtil.GetDb(), ctx, registerUserRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if numberOfUser > 0 {
		err = exception.NewBadRequestException("email already exists")
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	var user modelentity.User
	user.Username = sql.NullString{Valid: true, String: registerUserRequest.Username}
	user.Email = sql.NullString{Valid: true, String: registerUserRequest.Email}
	user.Password = sql.NullString{Valid: true, String: string(hashedPassword)}
	user.Utc = sql.NullString{Valid: true, String: registerUserRequest.Utc}
	user.CreatedAt = sql.NullInt64{Valid: true, Int64: time.Now().UnixMilli()}
	rowsAffected, err := service.UserRepository.Create(service.MysqlUtil.GetDb(), ctx, user)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not 1")
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	registerUserResponse = modelresponse.ToRegisterUserResponse(user)
	return
}

func (service *UserServiceImplementation) Login(ctx context.Context, requestId string, sessionIdUser string, loginUserRequest modelrequest.LoginUserRequest) (sessionId string, err error) {
	if sessionIdUser != "" {
		var rowsAffected int64
		rowsAffected, err = service.RedisUtil.GetClient().Del(ctx, sessionIdUser).Result()
		if err != nil && err != redis.Nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}
		if rowsAffected != 1 {
			err = errors.New("rows affected not 1")
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}
	}
	err = service.Validate.Struct(loginUserRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, loginUserRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}

	user, err := service.UserRepository.FindByEmail(service.MysqlUtil.GetDb(), ctx, loginUserRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	} else if err != nil && err == sql.ErrNoRows {
		err = exception.NewBadRequestException("wrong email or password")
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginUserRequest.Password))
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException("wrong email or password")
		return
	}

	userPermissions, err := service.UserPermissionRepository.FindByUserId(service.MysqlUtil.GetDb(), ctx, user.Id.Int32)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	var idPermissions []int32
	for _, userPermission := range userPermissions {
		idPermissions = append(idPermissions, userPermission.PermissionId.Int32)
	}

	sessionId = uuid.New().String()
	sessionValue := make(map[string]interface{})
	sessionValue["id"] = user.Id.Int32
	sessionValue["username"] = user.Username.String
	sessionValue["email"] = user.Email.String
	sessionValue["idPermissions"] = idPermissions
	sessionByte, err := json.Marshal(sessionValue)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	session := string(sessionByte)
	_, err = service.RedisUtil.GetClient().Set(ctx, sessionId, session, 0).Result()
	if err != nil && err != redis.Nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	return
}

func (service *UserServiceImplementation) Logout(ctx context.Context, requestId string, sessionId string) (err error) {
	rowsAffected, err := service.RedisUtil.GetClient().Del(ctx, sessionId).Result()
	if err != nil && err != redis.Nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not 1")
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	return
}
