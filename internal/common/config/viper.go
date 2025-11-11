package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewViperConfig() error {
	if err := godotenv.Load("../common/config/.env"); err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.AutomaticEnv()
	viper.BindEnv("stripe-key", "STRIPE_SECRET_KEY")
	viper.BindEnv("local-env", "LOCAL_ENV")

	return viper.ReadInConfig()
}
