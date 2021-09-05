package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DbConfig DbConfig `mapstructure:"db"`
}

type DbConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func ReadConfig() *AppConfig {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("BL")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// Note: Environment variables should be given full upper case

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return &config
}

func Dump() {
	fmt.Println("\nCONFIG")
	fmt.Println("===========")
	fmt.Printf("%+v\n", *ReadConfig())
}
