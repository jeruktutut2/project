package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	ProjectUserApplicationPort            string `mapstructure:"PROJECT_USER_APPLICATION_PORT"`
	ProjectUserApplicationTimeout         uint8  `mapstructure:"PROJECT_USER_APPLICATION_TIMEOUT"`
	ProjectUserMysqlHost                  string `mapstructure:"PROJECT_USER_MYSQL_HOST"`
	ProjectUserMysqlUsername              string `mapstructure:"PROJECT_USER_MYSQL_USERNAME"`
	ProjectUserMysqlPassword              string `mapstructure:"PROJECT_USER_MYSQL_PASSWORD"`
	ProjectUserMysqlDatabase              string `mapstructure:"PROJECT_USER_MYSQL_DATABASE"`
	ProjectUserMysqlMaxOpenConnection     uint16 `mapstructure:"PROJECT_USER_MYSQL_MAX_OPEN_CONNECTION"`
	ProjectUserMysqlMaxIdleConnection     uint16 `mapstructure:"PROJECT_USER_MYSQL_MAX_IDLE_CONNECTION"`
	ProjectUserMysqlConnectionMaxLifetime uint16 `mapstructure:"PROJECT_USER_MYSQL_CONNECTION_MAX_LIFETIME"`
	ProjectUserMysqlConnectionMaxIdletime uint16 `mapstructure:"PROJECT_USER_MYSQL_CONNECTION_MAX_IDLETIME"`
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
	println(time.Now().String() + " reading config file")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic("read config error: " + err.Error())
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic("unmarshal config error: " + err.Error())
	}
	println(time.Now().String() + " config file is read")
	return configuration
}
