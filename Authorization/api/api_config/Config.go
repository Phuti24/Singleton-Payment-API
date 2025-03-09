package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `yaml:"api_server"`
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type SecurityConfig struct {
	BcryptCost int `yaml:"bcrypt_cost"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // Config file name (without extension)
	viper.AddConfigPath("../")    // Look for the config file in the current directory
	viper.SetConfigType("yaml")   // The file type is YAML

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config into struct: %v", err)
	}

	return &config, nil
}
