package util

import "github.com/spf13/viper"

// Config stores all configurations of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	BitlyOauthLogin string `mapstructure:"BITLY_OAUTH_LOGIN"`
	BitlyOauthToken string `mapstructure:"BITLY_OAUTH_TOKEN"`
	Port            string `mapstructure:"PORT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("challenge")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
