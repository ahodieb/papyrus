package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ahodieb/papyrus/editor"
	"github.com/ahodieb/papyrus/notes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	APP_NAME          = "papyrus"
	CONFIG_ENV_PREFIX = "PAPYRUS"

	F_JOURNAL_FILE = "journal-file"
	F_CONFIG       = "config"
	F_EDITOR       = "editor"
)

var (
	cfgFile string
	version = "development"
)

var rootCmd = &cobra.Command{
	Use:     APP_NAME,
	Short:   "Tools to automate note taking workflow",
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		m := NewManager(cmd, args)
		return m.OpenLatest()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP(F_CONFIG, "c", "", "config file (default is $HOME/.papyrus.yaml)")
	rootCmd.PersistentFlags().StringP(F_JOURNAL_FILE, "j", "", "Journal file")
	rootCmd.PersistentFlags().StringP(F_EDITOR, "e", "", "Editor to use")
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
		viper.SetConfigName("." + APP_NAME)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(CONFIG_ENV_PREFIX)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	bindFlags()
}

func NewManager(cmd *cobra.Command, args []string) notes.Manager {
	file := getJournalFile(cmd)
	// FIXME do some default behavior
	if file == "" {
		cobra.CheckErr(fmt.Errorf("%q not specified", F_JOURNAL_FILE))
	}
	editorName := getEditor(cmd)

	editorName, err := cmd.Flags().GetString(F_EDITOR)
	cobra.CheckErr(err)

	return notes.Manager{
		Editor: editor.ByName(editorName),
		Notes:  notes.ReadFromFile(file),
	}
}

// FIXME: This only works for flags on root command, not sub commands
// This is most likely because this gets evaluated before the other commands are initialized or added as sub commands
// for the root command
// So far this is not a big deal because all three flags could be root flags (for now),
//  but i want to get this working generally to included it in my cli tool template
func bindFlags() {
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", CONFIG_ENV_PREFIX, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func getStringFlag(f string, cmd *cobra.Command) string {
	return cmd.Flags().Lookup(f).Value.String()
}

func getEditor(cmd *cobra.Command) string {
	return getStringFlag(F_EDITOR, cmd)
}

func getJournalFile(cmd *cobra.Command) string {
	return getStringFlag(F_JOURNAL_FILE, cmd)
}
