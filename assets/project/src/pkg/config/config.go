package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init 初始化viper
func New(path string) (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)

	v.SetConfigName(string(path))
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc//")  // path to look for the config file in
	v.AddConfigPath("$HOME/.") // call multiple times to add many search paths

	// v.AutomaticEnv()
	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		return nil, err
	}

	return v, err
}
