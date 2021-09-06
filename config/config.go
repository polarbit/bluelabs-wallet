package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Db       string `mapstructure:"db"`
	LogLevel string `mapstructure:"loglevel"`
}

var Config *AppConfig

func Init() {
	initConfig()
	initLogger()
}

func initConfig() {
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

	if _, err := zerolog.ParseLevel(config.LogLevel); err != nil {
		panic(fmt.Errorf("loglevel is incorrect, valid values are: panic, fatal, error, warning, debug, trace. err:%v", err))
	}

	Config = &config
}

func initLogger() {
	if level, err := zerolog.ParseLevel(Config.LogLevel); err != nil {
		panic("invalid log level in config:" + Config.LogLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}

	log.Logger.With().Logger()
}

func Dump() {
	fmt.Println("\nCONFIG")
	fmt.Println("===========")
	fmt.Printf("%+v\n", *Config)
}
