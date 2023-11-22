package utils

import "github.com/spf13/viper"

type Envs struct {
	DBsource      string `mapstructure:"DB_SOURCE"`
	DBdriver      string `mapstructure:"DB_DRIVER"`
	Address       string `mapstructure:"SERVER_ADDRESS"`
	JwtSecreteKey string `mapstructure:"JWT_SECRETE_KEY"`
}

func Load(path string) (*Envs, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var envs *Envs
	if err := viper.Unmarshal(&envs); err != nil {
		return nil, err
	}
	return envs, nil
}
