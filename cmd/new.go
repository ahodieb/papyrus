package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "add a new section",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := NewManager(cmd, args)
		p := m.AddEntry(strings.Join(args, " "), time.Now())
        
		return m.Open(p)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
