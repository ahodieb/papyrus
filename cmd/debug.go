package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Dump debug details",
	Long:  "Display different details useful for debugging",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Viper Config: %q\n", viper.ConfigFileUsed())
		if viper.ConfigFileUsed() == "" {
			fmt.Println("No config file specified")
		} else {
			viper.Debug()
		}
		fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>\n")

		fmt.Println("Journal File:", cmd.Flags().Lookup(F_JOURNAL_FILE).Value.String())
		fmt.Println("Editor:", cmd.Flags().Lookup(F_EDITOR).Value.String())
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
