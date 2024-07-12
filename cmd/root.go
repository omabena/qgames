package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qgames-parer",
		Short: "CLI Interface for parsing qgames logs",
	}

	cmd.AddCommand(NewCmd())
	return cmd
}
