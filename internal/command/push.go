package command

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mlouage/envtamer-go/internal/storage"
	"github.com/mlouage/envtamer-go/internal/util"
)

func newPushCmd() *cobra.Command {
	var filename string

	cmd := &cobra.Command{
		Use:   "push [DIRECTORY_NAME]",
		Short: "Push the contents of a local .env file to the database",
		Long:  `This command reads the specified .env file and stores its contents in the database, associated with the given directory.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve directory path
			var dirPath string
			var err error
			if len(args) > 0 {
				dirPath, err = util.ResolvePath(args[0])
			} else {
				dirPath, err = util.ResolvePath("")
			}
			if err != nil {
				return fmt.Errorf("failed to resolve directory path: %w", err)
			}

			// Parse .env file
			envFilePath := filepath.Join(dirPath, filename)
			envVars, err := util.ParseEnvFile(envFilePath)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}

			// Save to database
			db, err := storage.New()
			if err != nil {
				return fmt.Errorf("failed to create storage: %w", err)
			}
			defer db.Close()

			if err := db.SaveEnvVars(dirPath, envVars); err != nil {
				return fmt.Errorf("failed to save env vars: %w", err)
			}

			fmt.Printf("Successfully pushed %d environment variables for directory: %s\n", len(envVars), dirPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", ".env", "The name of the env file")
	return cmd
}
