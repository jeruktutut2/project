package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	ApplicationPort    string `mapstructure:"APPLICATION_PORT"`
	ApplicationTimeout uint8  `mapstructure:"APPLICATION_TIMEOUT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisDatabase      int    `mapstructure:"REDIS_DATABASE"`
	GrpcUserHost       string `mapstructure:"GRPC_USER_HOST"`
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
