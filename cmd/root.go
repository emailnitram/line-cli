package cmd

import (
	"os"

	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "line",
	Short: "LINE API CLI",
	Long:  "A CLI for the LINE Messaging API and LINE Login.",
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(messageCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(richMenuCmd)
	rootCmd.AddCommand(webhookCmd)
	rootCmd.AddCommand(audienceCmd)
	rootCmd.AddCommand(couponCmd)
rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(liffCmd)
}
