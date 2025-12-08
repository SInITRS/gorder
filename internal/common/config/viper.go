package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// NewViperConfig initializes the viper configuration.
func NewViperConfig() error {
	if err := godotenv.Load("../common/config/.env"); err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.AutomaticEnv()
	viper.BindEnv("stripe-key", "STRIPE_SECRET_KEY")
	viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	viper.BindEnv("local-env", "LOCAL_ENV")

	return viper.ReadInConfig()
}
