package grpc_test

import (
	"context"
	"log"
	"net"
	"project-user/helpers"
	"project-user/interceptor"
	"project-user/setup"
	"project-user/tests/initialize"
	"project-user/utils"
	"testing"
	"time"

	grpcuser "project-user/grpc"

	repositories "project-user/repositories"

	pbuser "project-user/grpc/pb/api/v1/user"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx                 context.Context
	requestId           string
	sessionId           string
	registerUserRequest *pbuser.RegisterRequest
	loginUserRequest    *pbuser.LoginRequest
	logoutUserRequest   *pbuser.LogoutRequest
	bufSize             int
	listener            *bufconn.Listener
	mysqlUtil           utils.MysqlUtil
	redisUtil           utils.RedisUtil
	validate            *validator.Validate
	bcryptHelper        helpers.BcryptHelper
	timeHelper          helpers.TimeHelper
	redisRepository     repositories.RedisRepository
	uuidHelper          helpers.UuidHelper
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (sut *UserServiceTestSuite) SetupSuite() {
	sut.T().Log("SetupSuite")
	sut.ctx = context.Background()
	sut.requestId = "requestId"
	sut.sessionId = "sessionId"
	sut.bufSize = 1024 * 1024
	sut.listener = bufconn.Listen(sut.bufSize)
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024), // Set max receive message size (1MB in this example)
		grpc.UnaryInterceptor(interceptor.SetLog),
	}
	grpcServer := grpc.NewServer(opts...)

	sut.mysqlUtil = utils.NewMysqlConnection("root", "12345", "localhost:3309", "user", 10, 10, 10, 10)
	sut.redisUtil = utils.NewRedisConnection("localhost", "6380", 0)
	sut.validate = validator.New()
	helpers.UsernameValidator(sut.validate)
	helpers.PasswordValidator(sut.validate)
	helpers.TelephoneValidator(sut.validate)
	sut.bcryptHelper = helpers.NewBcryptHelper()
	sut.timeHelper = helpers.NewTimeHelper()
	sut.redisRepository = repositories.NewRedisRepository()
	sut.uuidHelper = helpers.NewUuidHelper()

	repositorySetup := setup.NewRepositorySetup()
	serviceSetup := setup.NewServiceSetup(sut.mysqlUtil, sut.redisUtil, sut.validate, repositorySetup, sut.bcryptHelper, sut.timeHelper, sut.redisRepository, sut.uuidHelper)
	userGrpcService := grpcuser.NewUserGrpcService(serviceSetup.UserService)
	pbuser.RegisterUserServiceServer(grpcServer, userGrpcService)
	go func() {
		err := grpcServer.Serve(sut.listener)
		if err != nil {
			log.Fatalln(time.Now().String(), "error when server grpc:", err)
		}
	}()
}

func (sut *UserServiceTestSuite) SetupTest() {
	sut.T().Log("SetupTest")
	sut.registerUserRequest = &pbuser.RegisterRequest{
		Username:        "username",
		Email:           "email@email.com",
		Password:        "password@A1",
		Confirmpassword: "password@A1",
		Utc:             "+0800",
	}
	sut.loginUserRequest = &pbuser.LoginRequest{
		Email:    "email@email.com",
		Password: "password@A1",
	}
	sut.logoutUserRequest = &pbuser.LogoutRequest{
		Sessionid: "sessionId",
	}
}

func (sut *UserServiceTestSuite) BeforeTest(suiteName, testName string) {
	sut.T().Log("BeforeTest: " + suiteName + " " + testName)
}

func (sut *UserServiceTestSuite) TestRegisterRegisterUserRequestValidationError() {
	sut.T().Log("TestRegisterRegisterUserRequestValidationError")

	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	registerUserRequest := &pbuser.RegisterRequest{
		Username:        "",
		Email:           "",
		Password:        "",
		Confirmpassword: "",
		Utc:             "",
	}
	response, err := client.Register(sut.ctx, registerUserRequest)

	var registerResponse *pbuser.RegisterResponse
	sut.Equal(response, registerResponse)
	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"username","message":"is required"},{"field":"email","message":"is required"},{"field":"password","message":"is required"},{"field":"confirmpassword","message":"is required"},{"field":"utc","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterPasswordAndConfirmpasswordIsDifferent() {
	sut.T().Log("TestRegisterPasswordAndConfirmpasswordIsDifferent")
	sut.registerUserRequest.Confirmpassword = "password@A1-"

	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Register(sut.ctx, sut.registerUserRequest)

	var registerResponse *pbuser.RegisterResponse
	sut.Equal(response, registerResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"message","message":"password and confirm password is different"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameInternalServerError() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Register(sut.ctx, sut.registerUserRequest)

	var registerResponse *pbuser.RegisterResponse
	sut.Equal(response, registerResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.Internal)
	sut.Equal(st.Message(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByUsernameBadRequestUsernameAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByUsernameBadRequestUsernameAlreadyExists")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Register(sut.ctx, sut.registerUserRequest)

	var registerResponse *pbuser.RegisterResponse
	sut.Equal(response, registerResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"message","message":"username already exists"}]`)
}

func (sut *UserServiceTestSuite) TestRegisterUserRepositoryCountByCountByEmailBadRequestUsernameAlreadyExists() {
	sut.T().Log("TestRegisterUserRepositoryCountByCountByEmailBadRequestUsernameAlreadyExists")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)

	sut.registerUserRequest.Username = "username1"
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Register(sut.ctx, sut.registerUserRequest)

	var registerResponse *pbuser.RegisterResponse
	sut.Equal(response, registerResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"message","message":"email already exists"}]`)
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
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Register(sut.ctx, sut.registerUserRequest)
	sut.NotEqual(response, nil)
	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLoginRedisRepositoryDelWithSessionIdInternalServerError() {
	sut.T().Log("TestLoginRedisRepositoryDelWithSessionIdInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.loginUserRequest = &pbuser.LoginRequest{}
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
		"sessionid", sut.sessionId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Login(sut.ctx, sut.loginUserRequest)

	var loginResponse *pbuser.LoginResponse
	sut.Equal(response, loginResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"email","message":"is required"},{"field":"password","message":"is required"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("TestLoginUserRepositoryFindByEmailBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
		"sessionid", sut.sessionId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Login(sut.ctx, sut.loginUserRequest)

	var loginResponse *pbuser.LoginResponse
	sut.Equal(response, loginResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError() {
	sut.T().Log("TestLoginBcryptCompareHashAndPasswordBadRequestWrongEmailOrPasswordError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	sut.loginUserRequest.Password = "password@A1-"
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
		"sessionid", sut.sessionId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Login(sut.ctx, sut.loginUserRequest)

	var loginResponse *pbuser.LoginResponse
	sut.Equal(response, loginResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.InvalidArgument)
	sut.Equal(st.Message(), `[{"field":"message","message":"wrong email or password"}]`)
}

func (sut *UserServiceTestSuite) TestLoginUserPermissionRepositoryFindByUserIdInternalServerError() {
	sut.T().Log("TestLoginUserPermissionRepositoryFindByUserIdInternalServerError")
	initialize.DropTableUserPermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTablePermission(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.DropTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateTableUser(sut.mysqlUtil.GetDb(), sut.ctx)
	initialize.CreateDataUser(sut.mysqlUtil.GetDb(), sut.ctx)
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
		"sessionid", sut.sessionId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Login(sut.ctx, sut.loginUserRequest)

	var loginResponse *pbuser.LoginResponse
	sut.Equal(response, loginResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.Internal)
	sut.Equal(st.Message(), `[{"field":"message","message":"internal server error"}]`)
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
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
		"sessionid", sut.sessionId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Login(sut.ctx, sut.loginUserRequest)

	var loginResponse pbuser.LoginResponse
	loginResponse.Sessionid = ""
	sut.NotEqual(response, &loginResponse)

	sut.Equal(err, nil)
}

func (sut *UserServiceTestSuite) TestLogoutRowsAffectedNotOneInternalServerError() {
	sut.T().Log("TestLogoutRowsAffectedNotOneInternalServerError")
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Logout(sut.ctx, sut.logoutUserRequest)

	var logoutUserResponse *pbuser.LogoutResponse
	sut.Equal(response, logoutUserResponse)

	st, ok := status.FromError(err)
	if !ok {
		log.Fatalln("error when checking error status")
	}
	sut.Equal(st.Code(), codes.Internal)
	sut.Equal(st.Message(), `[{"field":"message","message":"internal server error"}]`)
}

func (sut *UserServiceTestSuite) TestLogoutSuccess() {
	sut.T().Log("TestLogoutSuccess")
	initialize.DelDataRedis(sut.redisUtil.GetClient(), sut.ctx, sut.sessionId)
	initialize.SetDataRedis(sut.redisUtil.GetClient(), sut.ctx, sut.sessionId, "value", time.Duration(0))
	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return sut.listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("error when creating new client:", err.Error())
	}
	defer conn.Close()

	client := pbuser.NewUserServiceClient(conn)
	md := metadata.Pairs(
		"requestid", sut.requestId,
	)
	sut.ctx = metadata.NewOutgoingContext(sut.ctx, md)
	response, err := client.Logout(sut.ctx, sut.logoutUserRequest)

	var logoutUserResponse pbuser.LogoutResponse
	logoutUserResponse.Msg = ""
	sut.NotEqual(response, &logoutUserResponse)

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
