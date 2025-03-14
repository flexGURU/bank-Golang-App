package utils

import (
	"time"

	"github.com/spf13/viper"
)


type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DSN string `mapstructure:"DSN"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	From_Email string `mapstructure:"FROM_EMAIL"`
	EmailSendName string `mapstructure:"EMAIL_SEND_NAME"`
	EamilPassword string `mapstructure:"EMAIL_PASSWORD"`
	Origin string `mapstructure:"ORIGIN"`
}


func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return 
	}

	err = viper.Unmarshal(&config)
	return

	
}


