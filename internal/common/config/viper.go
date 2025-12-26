package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() (err error) {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}

	// Load .env into os environment
	loadEnvToOs(relPath)

	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.AutomaticEnv()
	viper.BindEnv("stripe-key", "STRIPE_SECRET_KEY")
	viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	viper.BindEnv("local-env", "LOCAL_ENV")

	return viper.ReadInConfig()
}

func loadEnvToOs(path string) {
	v := viper.New()
	v.SetConfigFile(filepath.Join(path, ".env"))
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err == nil {
		for _, key := range v.AllKeys() {
			os.Setenv(strings.ToUpper(key), v.GetString(key))
		}
	} else {
		fmt.Printf("Failed to load .env: %v\n", err)
	}
}

func getRelativePathFromCaller() (relPath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, here, _, _ := runtime.Caller(0)
	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	fmt.Printf("caller from: '%s', here: '%s', relpath: '%s'\n", callerPwd, here, relPath)
	return
}
