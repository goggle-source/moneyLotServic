package config

import (
	"time"

	"github.com/spf13/viper"
)

type Cfg struct {
	Env  string   `mapstructure:"Env"`
	Path string   `mapstructure:"Path"`
	GRPC GRPC     `mapstructure:"GRPC"`
	DB   Database `mapstructure:"Database"`
}

type GRPC struct {
	Timeout time.Duration `mapstructure:"timeout"`
	Port    int           `mapstructure:"port"`
}

type Database struct {
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DbName       string        `mapstructure:"dbName"`
	Port         int           `mapstructure:"port"`
	HostDB       string        `mapstructure:"hostDB"`
	MaxOpenConns int           `mapstructure:"maxOpenConns"`
	MaxIdleConns int           `mapstructure:"maxIndleConns"`
	ConnMaxIdle  time.Duration `mapstructure:"connMaxIdle"`
	ConnLifeTime time.Duration `mapstructure:"connLifeTime"`
}

func MustLoad() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	var cfg Cfg
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
