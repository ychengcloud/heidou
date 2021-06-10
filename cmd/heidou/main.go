package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "Heidou",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		for i, arg := range args {
			fmt.Printf("arg %d : %s\n", i, arg)
		}
		cmd.Help()
	},
}

var versionTemplate = `
	{{with .Name}}{{printf "%s: " .}}{{end}}
	{{printf "version %s" .Version}}
	Commit: ` + commit + `
	Build Date: ` + date + `

`

func init() {
	rootCmd.SetVersionTemplate(versionTemplate)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
