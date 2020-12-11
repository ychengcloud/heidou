package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/horcus/heidou"
)

var configFilename string
var pkgPath string

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

		if err := genConfig(filepath.Join(projectName, configFilename), pkgPath); err != nil {
			return
		}

		err := heidou.GenProject(projectName, pkgPath)
		if err != nil {
			fmt.Println("err:", err)
		}
		fmt.Println("Done")
	},
}

func init() {
	initCmd.Flags().StringVarP(&configFilename, "config", "c", "heidou-example.yml", "config file name")
	initCmd.Flags().StringVarP(&pkgPath, "pkg-name", "p", "", "package name")
	err := initCmd.MarkFlagRequired("pkg-name")
	if err != nil {
		return
	}
	rootCmd.AddCommand(initCmd)

}

func genConfig(name string, pkgPath string) error {
	if name == "" {
		name = "heidou-example.yml"
	}

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return fmt.Errorf("unable to create config dir: " + err.Error())
	}

	var buf bytes.Buffer
	if err := configTemplate.Execute(&buf, pkgPath); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(name, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("unable to write cfg file: " + err.Error())
	}

	return nil
}

var configTemplate = template.Must(template.New("name").Parse(
	`pkgPath: "{{.}}"

db:
  dialect: mysql
  user: dbuser
  password: dbpass
  host: "127.0.0.1"
  port: 3306
  name: dbname
  charset: utf8mb4

tables:
  - name: product
    fields:
    - name: id
      isRequired: true
      isFilterable: true
    - name: name
      isRequired: true
      isFilterable: true
      tags: 'binding:"required,max=64"'
    - name: variant
      joinType: "HasMany"
      tableName: product_variant
    - name: category
      joinType: "ManyToMany"
      tableName: category  
      joinTableName: product_category_relation
  - name: employee
    methods: ["list", "update", "delete", "bulkGet", "bulkDelete"]
    fields:
    - name: role
      joinType: "ManyToMany"
      tableName: role  
      joinTableName: employee_role_relation
  - name: role
    methods: ["list", "get", "bulkGet", "bulkDelete"]
    fields:
    - name: action
      joinType: "ManyToMany"
      tableName: action  
      joinTableName: role_action_relation
      errorCodes: [E1, E2]
    - name: resource
      joinType: "ManyToMany"
      tableName: resource  
      joinTableName: role_resource_relation
      errorCodes: [E1, E2]
  - name: casbin_rules
    isSkip: true
  - name: employee_role_relation
    isSkip: true
  - name: migrations
    isSkip: true
  - name: product_attribute_value_relation
    isSkip: true
  - name: product_category_relation
    isSkip: true
  - name: product_collection_relation
    isSkip: true
  - name: product_midea_file_relation
    isSkip: true
  - name: product_type_variant_attribute_relation
    isSkip: true
  - name: product_variant_attribute_value_relation
    isSkip: true
  - name: role_action_relation
    isSkip: true
  - name: role_resource_relation
    isSkip: true
  - name: event_logs
    isSkip: true
  
`))
