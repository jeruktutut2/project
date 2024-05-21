package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	ApplicationPort            string `mapstructure:"APPLICATION_PORT"`
	ApplicationTimeout         uint8  `mapstructure:"APPLICATION_TIMEOUT"`
	MysqlHost                  string `mapstructure:"MYSQL_HOST"`
	MysqlUsername              string `mapstructure:"MYSQL_USERNAME"`
	MysqlPassword              string `mapstructure:"MYSQL_PASSWORD"`
	MysqlDatabase              string `mapstructure:"MYSQL_DATABASE"`
	MysqlMaxOpenConnection     uint16 `mapstructure:"MYSQL_MAX_OPEN_CONNECTION"`
	MysqlMaxIdleConnection     uint16 `mapstructure:"MYSQL_MAX_IDLE_CONNECTION"`
	MysqlConnectionMaxLifetime uint16 `mapstructure:"MYSQL_CONNECTION_MAX_LIFETIME"`
	MysqlConnectionMaxIdletime uint16 `mapstructure:"MYSQL_CONNECTION_MAX_IDLETIME"`
	RedisHost                  string `mapstructure:"REDIS_HOST"`
	RedisDatabase              int    `mapstructure:"REDIS_DATABASE"`
	JwtKey                     string `mapstructure:"JWT_KEY"`
	JwtAccessTokenExpireTime   uint16 `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE_TIME"`
	JwtRefreshTokenExpireTime  uint16 `mapstructure:"JWT_REFRESH_TOKEN_EXPIRE_TIME"`
	RabbitmqHost               string `mapstructure:"RABBITMQ_HOST"`
	RabbitmqUsername           string `mapstructure:"RABBITMQ_USERNAME"`
	RabbitmqPassword           string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitmqPort               string `mapstructure:"RABBITMQ_PORT"`
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
