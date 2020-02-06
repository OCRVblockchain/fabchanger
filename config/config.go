package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	General
	Mode string
}

type General struct {
	Identity          string
	MSPId             string
	Channel           string
	ConnectionProfile string
	ConfigTxPath      string
	OrgToJoinMSP      string
	MyOrg             string
}

func GetConfig() (*Config, error) {

	var Configuration *Config
	mode := flag.String("mode", "", "operating mode")
	flag.Parse()

	if *mode == "" {
		return nil, errors.New("please specify --mode")
	}

	// read config
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read config file, %s", err))
	}

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to decode into struct, %v", err))
	}
	Configuration.Mode = *mode
	return Configuration, nil
}
