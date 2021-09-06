package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Db string `mapstructure:"db"`
}

func ReadConfig() *AppConfig {

	// Config file
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(findRootDir())
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// TODO: Move this note to README.
	// Environment variables
	// Note: Environment variables should be given full upper case
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &config
}

func Dump() {
	fmt.Println("\nCONFIG")
	fmt.Println("===========")
	fmt.Printf("%+v\n", *ReadConfig())
}
