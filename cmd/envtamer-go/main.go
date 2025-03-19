package main

import (
	"fmt"
	"os"

	"github.com/mlouage/envtamer-go/internal/command"
)

func main() {
	rootCmd := command.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
