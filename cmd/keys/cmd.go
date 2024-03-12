package list

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:     "keys",
		Aliases: []string{"ls", "l"},
		Short:   "List user's SSH keys on GitHub",
		RunE:    run,
	}
}

func run(cmd *cobra.Command, args []string) (err error) {
	return nil
}
