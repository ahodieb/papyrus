package cmd

import (
	"fmt"
	"os"

	"github.com/ahodieb/papyrus/editor"
	"github.com/ahodieb/papyrus/notes"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	JOURNAL_FILE = "journal-file"
	EDITOR       = "editor"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "papyrus",
	Short: "Tools to automate note taking workflow",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.papyrus.yaml)")
	openCmd.PersistentFlags().String(EDITOR, "", "Editor to use")
	rootCmd.PersistentFlags().String(JOURNAL_FILE, "", "Journal file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)

		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(wd)
		viper.AddConfigPath(home)
		viper.SetConfigName(".papyrus")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PAPYRUS")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func NewManager(cmd *cobra.Command, args []string) (notes.Manager, error) {
	file, _ := cmd.Flags().GetString(JOURNAL_FILE)
	editorName, _ := cmd.Flags().GetString(EDITOR)
	n, err := notes.ReadOrCreate(file)
	if err != nil {
		return notes.Manager{}, err
	}

	return notes.Manager{
		Editor: editor.ByName(editorName),
		Notes:  n,
	}, nil
}
