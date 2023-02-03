package config

import (
	"Tiktok/pkg/log"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type Jaeger struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

var ServerSetting = &Server{}
var DatabaseSetting = &Database{}
var JaegerSetting = &Jaeger{}

func InitViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		log.Error("Read config file err", zap.Error(err))
	}
	log.Info("Config init success")

	server := viper.Sub("server")
	if err = server.Unmarshal(&ServerSetting); err != nil {
		log.Error("Config load err", zap.Error(err))
	}

	database := viper.Sub("database")
	if err = database.Unmarshal(&DatabaseSetting); err != nil {
		log.Error("Config load err", zap.Error(err))
	}

	jaeger := viper.Sub("jaeger")
	if err = jaeger.Unmarshal(&JaegerSetting); err != nil {
		log.Error("Config load err", zap.Error(err))
	}
}
