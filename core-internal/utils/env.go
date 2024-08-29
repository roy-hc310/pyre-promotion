package utils

import (
	"pyre-promotion/core-internal/conf"

	"github.com/spf13/viper"
)

var GlobalEnv conf.Env

func LoadGlobalEnv(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.Unmarshal(&GlobalEnv)
	viper.Unmarshal(&GlobalEnv.DBRead)
	viper.Unmarshal(&GlobalEnv.DBWrite)
	viper.Unmarshal(&GlobalEnv.Kafka)
	viper.Unmarshal(&GlobalEnv.Redis)
	viper.Unmarshal(&GlobalEnv.Elastic)
	viper.Unmarshal(&GlobalEnv.KeyCloak)

	return nil
}
