package cmd

import (
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open journal file",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := NewManager(cmd, args)
		return m.OpenLatest()
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
