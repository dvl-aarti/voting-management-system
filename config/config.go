package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	JWT_SECRET_KEY       string        `mapstructure:"JWT_SECRET_KEY"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	TestDatabase         string        `mapstructure:"TEST_DB"`
	Database             string        `mapstructure:"DATABASE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

var APP_CONFIG config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to read config")
	}

	_ = viper.Unmarshal(&APP_CONFIG)
}
