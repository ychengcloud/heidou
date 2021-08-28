package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
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

		cfg.TemplatesPath = tplPath

		if cfg.Overwrite {
			prompt := &survey.Confirm{
				Message: `[Warning]
The overwrite flag (Overwrite is True) is specified in the configuration file. If Yes is selected, the generated file will overwrite the existing file. 
Are you sure to continue?
				
配置文件中指定了覆写标志(Overwrite is True)，如果选择 Yes， 生成的文件将覆盖已有文件。
确认继续吗?`,
			}

			overwrite := false
			// ask the question
			err := survey.AskOne(prompt, &overwrite)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if err != nil || !overwrite {
				return
			}
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

	generateCmd.MarkFlagRequired("templates")
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
