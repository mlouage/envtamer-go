package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mlouage/envtamer-go/internal/storage"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize an empty database in the user's home folder",
		Long:  `This command creates an empty SQLite database file named 'envtamer-go.db' in the '.envtamer-go' directory of the user's home folder. If the file already exists, the command will not overwrite it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.New()
			if err != nil {
				return fmt.Errorf("failed to create storage: %w", err)
			}
			defer db.Close()

			if err := db.Init(); err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}

			fmt.Println("Database initialized successfully.")
			return nil
		},
	}

	return cmd
}
