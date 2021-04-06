package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ychengcloud/heidou"
)

var configFilename string

var initCmd = &cobra.Command{
	Use:     "init [flags] NAME",
	Short:   `initialize framework for project`,
	Example: `heidou init example`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("project name must be specified")
		}

		projectName := args[0]
		_, err := os.Stat(projectName)
		if !os.IsNotExist(err) {
			return errors.New("target path exist")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		if err := genConfig(projectName, configFilename); err != nil {
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	initCmd.Flags().StringVarP(&configFilename, "config", "c", "heidou-example.yml", "config file name")
	// initCmd.Flags().StringVarP(&pkgPath, "pkg-name", "p", "", "package name")
	// err := initCmd.MarkFlagRequired("pkg-name")
	// if err != nil {
	// 	return
	// }
	rootCmd.AddCommand(initCmd)

}

func genConfig(projectName, cfgName string) error {
	if cfgName == "" {
		cfgName = "heidou-example.yml"
	}

	name := filepath.Join(projectName, configFilename)

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return fmt.Errorf("unable to create config dir: " + err.Error())
	}

	path := "_assets/templates/heidou-example.yml.tmpl"
	tmpl := template.New(filepath.Base(path))

	configTemplate := template.Must(tmpl.ParseFS(heidou.Assets, path))

	var buf bytes.Buffer
	if err := configTemplate.Execute(&buf, projectName); err != nil {
		panic(err)
	}

	if err := os.WriteFile(name, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("unable to write cfg file: " + err.Error())
	}

	return nil
}
