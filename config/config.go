package config

import (
	"github.com/spf13/viper"
)

type apiConfig struct {
	Server struct {
		Port string `yaml:"port"`
	}

	DataFile string `yaml:"data_file"`
}

var Api apiConfig

func LoadConfig(env string) error {
	name := "config." + env
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Api)
	if err != nil {
		return err
	}

	return nil
}
