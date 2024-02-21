package util

import (
	"time"

	"github.com/spf13/viper"
)

const ConfigName = "config"
const ConfigType = "yaml"

var Configuration Config

type Config struct {
	Server struct {
		Mode            string        `mapstructure:"mode"`
		Port            int           `mapstructure:"port"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
		Endpoint        string        `mapstructure:"endpoint"`
	} `mapstructure:"server"`

	Platform struct {
		Host   string `mapstructure:"host"`
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"platform"`

	MongoDB struct {
		Host       string   `mapstructure:"host"`
		Port       int      `mapstructure:"port"`
		Database   string   `mapstructure:"database"`
		Collection string   `mapstructure:"collection"`
		Username   string   `mapstructure:"username"`
		Password   string   `mapstructure:"password"`
		Options    []string `mapstructure:"options"`
	} `mapstructure:"mongo"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	var config Config
	err = viper.Unmarshal(&config)
	Configuration = config
	return
}
