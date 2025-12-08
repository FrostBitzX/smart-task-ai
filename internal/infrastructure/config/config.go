package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	FrontendURL string `mapstructure:"FRONTEND_URL"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBName      string `mapstructure:"DB_DATABASE"`
	DBUsername  string `mapstructure:"DB_USERNAME"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBPort      string `mapstructure:"DB_PORT"`
}

func NewConfig() *Config {
	config := &Config{}

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Unable to read .env file", err)
		}
	}

	viper.AutomaticEnv()

	viper.BindEnv("FRONTEND_URL")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_DATABASE")
	viper.BindEnv("DB_USERNAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_PORT")

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return config
}
