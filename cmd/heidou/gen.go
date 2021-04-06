package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ychengcloud/heidou"
	"github.com/ychengcloud/heidou/api"
)

var configPath string
var tplPath string

var generateCmd = &cobra.Command{
	Use:     "generate [flags]",
	Short:   "generate go code for the database schema",
	Example: `heidou generate -c ./config.yml -t tpl_path`,
	Args: func(_ *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig(configPath)
		logrus.Println(cfg, cfg.DBConfig, cfg.Tables)
		for _, table := range cfg.Tables {
			logrus.Println(table)
		}
		if len(tplPath) > 0 {
			cfg.TemplatesPath = tplPath
		}
		err := api.Generate(cfg)
		if err != nil {
			log.Fatal("gen:", err)
		}
		fmt.Println("Done")
	},
}

func init() {
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "./config.yml", "config file path")
	generateCmd.Flags().StringVarP(&tplPath, "templates", "t", "", "template path")

	cobra.OnInitialize()
	rootCmd.AddCommand(generateCmd)

}

func loadConfig(path string) *heidou.Config {
	viper.SetConfigFile(string(path))

	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.AddConfigPath("/etc/heidou/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.heidou") // call multiple times to add many search paths
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		logrus.Fatalf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}

	cfg := &heidou.Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		logrus.Fatalf("Config unmarshal fail: %s \n", err)
		os.Exit(1)
	}
	return cfg
}
