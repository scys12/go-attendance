package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver   string `mapstructure:"db_driver"`
	DBName     string `mapstructure:"db_name"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
}

type ServerConfig struct {
	WriteTimeout int `mapstructure:"write_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	Port         int `mapstructure:"port"`
}

type JWTConfig struct {
	JWTSecret  string `mapstructure:"jwt_secret"`
	JWTExpires int    `mapstructure:"jwt_expires"`
}

func InitDBConfig() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	return
}

func InitServerConfig() (config ServerConfig, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	return
}

func InitJWTConfig() (config JWTConfig, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	return
}
