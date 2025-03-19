package command

import (
	"fmt"

	"github.com/spf13/cobra"
)
const logo = `

   ___   _ __   __   __ | |_    __ _   _ __ ___     ___   _ __
  / _ \ | '_ \  \ \ / / | __|  / _' | | '_ ' _ \   / _ \ | '__|
 |  __/ | | | |  \ V /  | |_  | (_| | | | | | | | |  __/ | |
  \___| |_| |_|   \_/    \__|  \__,_| |_| |_| |_|  \___| |_|

`

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envtamer",
		Short: "Taming digital environment files chaos with elegant simplicity",
		Long: fmt.Sprintf("%s\nA command-line tool for managing environment variables across different projects and directories.", logo),
	}

	// Add subcommands
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newPushCmd())
	cmd.AddCommand(newPullCmd())
	cmd.AddCommand(newListCmd())

	return cmd
}
