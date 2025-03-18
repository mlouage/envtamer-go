/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const logo = `

   ___   _ __   __   __ | |_    __ _   _ __ ___     ___   _ __
  / _ \ | '_ \  \ \ / / | __|  / _' | | '_ ' _ \   / _ \ | '__|
 |  __/ | | | |  \ V /  | |_  | (_| | | | | | | | |  __/ | |
  \___| |_| |_|   \_/    \__|  \__,_| |_| |_| |_|  \___| |_|

`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envtamer-go",
	Short: "Taming digital environment files chaos with elegant simplicity.",
	Long:  fmt.Sprintf("%s\nA command-line tool for managing environment variables across different projects and directories.", logo),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envtamer-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Long = fmt.Sprintf("%s\nA command-line tool for managing environment variables across different projects and directories.", logo)
}
