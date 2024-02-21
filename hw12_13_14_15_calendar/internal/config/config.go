package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	Host   HostConf
	DB     DBConf
}

type LoggerConf struct {
	Level string
}

type HostConf struct {
	Host string
	Port int
}

type DBConf struct {
	Type     string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func NewConfig(configPath string) Config {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return Config{
		LoggerConf{
			viper.GetString("logger.level"),
		},
		HostConf{
			viper.GetString("server.host"),
			viper.GetInt("server.port"),
		},
		DBConf{
			viper.GetString("storage.type"),
			viper.GetString("storage.host"),
			viper.GetInt("storage.port"),
			viper.GetString("storage.database"),
			viper.GetString("storage.username"),
			viper.GetString("storage.password"),
		},
	}
}
