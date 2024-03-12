package cmd

import (
	"github.com/helaili/gh-ssh-cert-please/cmd/keys"
	"github.com/spf13/cobra"
)

func New(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ssh-cert-please",
		Short:             "Manage SSH certificates for GitHub",
		PersistentPreRunE: preRun,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		Version:           buildVersion(version, commit),
	}

	cmd.AddCommand(
		keys.New(),
	)
	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	// cmd.SilenceUsage = true
	return nil
}

func buildVersion(version, commit string) string {
	if commit != "" {
		version += " (" + commit + ")"
	}
	return version
}
