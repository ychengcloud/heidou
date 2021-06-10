package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version contains the current version.
	Version = "dev"
	Commit  = ""
	// BuildDate contains a string with the build date.
	Date = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "Heidou",
	Version: Version,
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
	Commit: ` + Commit + `
	Build Date: ` + Date + `

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
