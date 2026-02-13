// gin-demo/internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey string
	Port   string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}

	AppConfig = Config{
		APIKey: viper.GetString("API_KEY"),
		Port:   viper.GetString("PORT"),
	}
}
