package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	PORT   string
	loglvl string
	fffff  string
	//dbhostname string       `mapstructure:"DBhostname"`
	//dbPort     int          `mapstructure:"DBPort"`
	//dbUser     string       `mapstructure:"DBUser"`
	//dbPassword string       `mapstructure:"DBPassword"`
	//dbname     string       `mapstructure:"DBname"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(config)
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(config)
	return
}
