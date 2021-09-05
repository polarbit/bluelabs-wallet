package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
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

// exists returns whether the given file or directory exists
func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// travers up directories until we find a main.go file
func findRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	for i := 0; i < 3; i++ {
		if ok, _ := isFileExists(path.Join(dir, "main.go")); ok {
			return dir
		}
		dir = path.Dir(dir)
	}

	panic(fmt.Errorf("fatal error config file: %w", err))
}
