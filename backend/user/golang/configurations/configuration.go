package configuration

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	ProjectUserApplicationHost            string `mapstructure:"PROJECT_USER_APPLICATION_HOST"`
	ProjectUserApplicationTimeout         uint8  `mapstructure:"PROJECT_USER_APPLICATION_TIMEOUT"`
	ProjectUserMysqlHost                  string `mapstructure:"PROJECT_USER_MYSQL_HOST"`
	ProjectUserMysqlUsername              string `mapstructure:"PROJECT_USER_MYSQL_USERNAME"`
	ProjectUserMysqlPassword              string `mapstructure:"PROJECT_USER_MYSQL_PASSWORD"`
	ProjectUserMysqlDatabase              string `mapstructure:"PROJECT_USER_MYSQL_DATABASE"`
	ProjectUserMysqlMaxOpenConnection     int    `mapstructure:"PROJECT_USER_MYSQL_MAX_OPEN_CONNECTION"`
	ProjectUserMysqlMaxIdleConnection     int    `mapstructure:"PROJECT_USER_MYSQL_MAX_IDLE_CONNECTION"`
	ProjectUserMysqlConnectionMaxLifetime int    `mapstructure:"PROJECT_USER_MYSQL_CONNECTION_MAX_LIFETIME"`
	ProjectUserMysqlConnectionMaxIdletime int    `mapstructure:"PROJECT_USER_MYSQL_CONNECTION_MAX_IDLETIME"`
	ProjectUserRedisHost                  string `mapstructure:"PROJECT_USER_REDIS_HOST"`
	ProjectUserRedisPort                  string `mapstructure:"PROJECT_USER_REDIS_PORT"`
	ProjectUserRedisDatabase              int    `mapstructure:"PROJECT_USER_REDIS_DATABASE"`
	JwtKey                                string `mapstructure:"JWT_KEY"`
	JwtAccessTokenExpireTime              uint16 `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE_TIME"`
	JwtRefreshTokenExpireTime             uint16 `mapstructure:"JWT_REFRESH_TOKEN_EXPIRE_TIME"`
	RabbitmqHost                          string `mapstructure:"RABBITMQ_HOST"`
	RabbitmqUsername                      string `mapstructure:"RABBITMQ_USERNAME"`
	RabbitmqPassword                      string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitmqPort                          string `mapstructure:"RABBITMQ_PORT"`
}

func NewConfiguration() (configuration *Configuration) {
	println(time.Now().String() + " reading environment variables")
	var conf Configuration
	conf.ProjectUserApplicationHost = os.Getenv("PROJECT_USER_APPLICATION_HOST")
	conf.ProjectUserMysqlHost = os.Getenv("PROJECT_USER_MYSQL_HOST")
	conf.ProjectUserMysqlUsername = os.Getenv("PROJECT_USER_MYSQL_USERNAME")
	conf.ProjectUserMysqlPassword = os.Getenv("PROJECT_USER_MYSQL_PASSWORD")
	conf.ProjectUserMysqlDatabase = os.Getenv("PROJECT_USER_MYSQL_DATABASE")
	var err error
	conf.ProjectUserMysqlMaxOpenConnection, err = strconv.Atoi(os.Getenv("PROJECT_USER_MYSQL_MAX_OPEN_CONNECTION"))
	if err != nil {
		log.Fatalln("error when convert mysql max open connection from string to int:", err)
	}
	conf.ProjectUserMysqlMaxIdleConnection, err = strconv.Atoi(os.Getenv("PROJECT_USER_MYSQL_MAX_IDLE_CONNECTION"))
	if err != nil {
		log.Fatalln("error when convert mysql max idle connection from string to int:", err)
	}
	conf.ProjectUserMysqlConnectionMaxLifetime, err = strconv.Atoi(os.Getenv("PROJECT_USER_MYSQL_CONNECTION_MAX_LIFETIME"))
	if err != nil {
		log.Fatalln("error when convert mysql connection max lifetime from string to int:", err)
	}
	conf.ProjectUserMysqlConnectionMaxIdletime, err = strconv.Atoi(os.Getenv("PROJECT_USER_MYSQL_CONNECTION_MAX_IDLETIME"))
	if err != nil {
		log.Fatalln("error when convert mysql connection max idle time from string to int:", err)
	}
	conf.ProjectUserRedisHost = os.Getenv("PROJECT_USER_REDIS_HOST")
	conf.ProjectUserRedisPort = os.Getenv("PROJECT_USER_REDIS_PORT")
	conf.ProjectUserRedisDatabase, err = strconv.Atoi(os.Getenv("PROJECT_USER_REDIS_DATABASE"))
	if err != nil {
		log.Fatalln("error when convert redis db from string to int:", err)
	}
	println(time.Now().String() + " environment variables is read")
	configuration = &conf
	return
}
