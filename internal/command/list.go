package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mlouage/envtamer-go/internal/storage"
	"github.com/mlouage/envtamer-go/internal/util"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [DIRECTORY_NAME]",
		Short: "List stored directories or environment variables",
		Long:  `If no directory is specified, this command lists all directories stored in the database. If a directory is provided, it lists all environment variables stored for that directory.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.New()
			if err != nil {
				return fmt.Errorf("failed to create storage: %w", err)
			}
			defer db.Close()

			if len(args) == 0 {
				// List directories
				directories, err := db.ListDirectories()
				if err != nil {
					return fmt.Errorf("failed to list directories: %w", err)
				}

				if len(directories) == 0 {
					fmt.Println("No directories stored in the database.")
					return nil
				}

				fmt.Println("Stored directories:")
				for _, dir := range directories {
					fmt.Println(dir)
				}
			} else {
				// List environment variables for a directory
				dirPath, err := util.ResolvePath(args[0])
				if err != nil {
					return fmt.Errorf("failed to resolve directory path: %w", err)
				}

				envVars, err := db.GetEnvVars(dirPath)
				if err != nil {
					return fmt.Errorf("failed to get env vars: %w", err)
				}

				if len(envVars) == 0 {
					fmt.Printf("No environment variables stored for directory: %s\n", dirPath)
					return nil
				}

				fmt.Printf("Environment variables for directory: %s\n", dirPath)
				for key, value := range envVars {
					fmt.Printf("%s=%s\n", key, value)
				}
			}

			return nil
		},
	}

	return cmd
}
