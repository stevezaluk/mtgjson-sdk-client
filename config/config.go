package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "/.config/mtgjson-cmd"
	defaultConfigName = "config.json"
)

/*
ReadConfigFile Parse the config file passed in the path to ensure that the SDK has the values we need
*/
func ReadConfigFile(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.SetConfigType("json")
		viper.AddConfigPath(home + defaultConfigPath)
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetDefault("api.use_ssl", false)
	viper.SetDefault("api.ip_address", "127.0.0.1")
	viper.SetDefault("api.port", 8080)
	viper.SetDefault("api.base_url", BuildBaseUrl(viper.GetString("api.ip_address"), viper.GetString("api.port"), viper.GetBool("api.use_ssl")))

	return nil
}

/*
BuildBaseUrl Build a base url string from an IP Address and a Port
*/
func BuildBaseUrl(ipAddress string, port string, useSsl bool) string {
	ret := "http://"
	if useSsl {
		ret = "https://"
	}

	ret += ipAddress + ":" + port + "/api/v1"

	return ret
}
