package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
