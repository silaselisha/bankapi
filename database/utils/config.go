package utils

import "github.com/spf13/viper"

type envs struct {
	DBsource string `mapstructure:"DB-SOURCE"`
	DBdriver string `mapstructure:"DB_DRIVER"`
	Address  string `mapstructure:"SERVER_ADDRESS"`
}

func Load(path string) (*envs, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var envs *envs
	if err := viper.Unmarshal(&envs); err != nil {
		return nil, err
	}
	return envs, nil
}
