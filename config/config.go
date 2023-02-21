package config

import (
	"Tiktok/pkg/log"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server 添加yaml字段便于后面直接从yaml类型的配置文件中关联绑定
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
	MaxConn  int    `yaml:"maxConn"`
	MaxOpen  int    `yaml:"maxOpen"`
}

type Jaeger struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

type Minio struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}

type Auth struct {
	Md5Salt   string `yaml:"md5Salt"`
	JwtSecret string `yaml:"jwtSecret"`
}

var ServerSetting = &Server{}
var DatabaseSetting = &Database{}
var JaegerSetting = &Jaeger{}
var AuthSetting = &Auth{}
var MinioSetting = &Minio{}

func InitViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		log.Error("Read config file err", zap.Error(err))
	}

	//获取viper的子结构（子树），否则字段差了一个层级没法绑定
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

	auth := viper.Sub("auth")
	if err = auth.Unmarshal(&AuthSetting); err != nil {
		log.Error("Config load err", zap.Error(err))
	}

	minio := viper.Sub("minio")
	if err = minio.Unmarshal(&MinioSetting); err != nil {
		log.Error("Config load err", zap.Error(err))
	}
	log.Info("Config init successfully")
}
