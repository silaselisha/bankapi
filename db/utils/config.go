package utils


import (
  "github.com/spf13/viper"
)

type Config struct {
   DBsource string `mapstructure:"DATA_SOURCE"`
   DBdriver string `mapstructure:"DB_DRIVER"`
   ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}


func LoadConfig(path string)(config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("path")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
	   return
	}

	err = viper.Unmarshal(&config)
	return
}
