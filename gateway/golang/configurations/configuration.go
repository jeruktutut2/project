package configuration

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	ProjectGatewayApplicationPort     string `mapstructure:"PROJECT_GATEWAY_APPLICATION_PORT"`
	ProjectGatewayApplicationTimeout  int    `mapstructure:"PROJECT_GATEWAY_APPLICATION_TIMEOUT"`
	ProjectGatewayUserApplicationHost string `mapstructure:"PROJECT_GATEWAY_USER_APPLICATION_HOST"`
	ProjectUserApplicationHost        string `mapstructure:"PROJECT_USER_APPLICATION_HOST"`
	ProjectUserApplicationPort        string `mapstructure:"PROJECT_USER_APPLICATION_PORT"`
}

func NewConfiguration() (configuration *Configuration) {
	println(time.Now().String() + " reading environment variables")
	var conf Configuration
	var err error
	conf.ProjectGatewayApplicationPort = os.Getenv("PROJECT_GATEWAY_APPLICATION_PORT")
	conf.ProjectGatewayApplicationTimeout, err = strconv.Atoi(os.Getenv("PROJECT_GATEWAY_APPLICATION_TIMEOUT"))
	if err != nil {
		log.Fatalln("error when convert project gateway application port from string to int:", err)
	}
	conf.ProjectGatewayUserApplicationHost = os.Getenv("PROJECT_GATEWAY_USER_APPLICATION_HOST")
	conf.ProjectUserApplicationHost = os.Getenv("PROJECT_USER_APPLICATION_HOST")
	println(time.Now().String() + " environment variables is read")
	configuration = &conf
	return configuration
}
