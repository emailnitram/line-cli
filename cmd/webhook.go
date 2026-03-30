package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage webhook endpoint",
}

var webhookGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get webhook endpoint information",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		data, err := c.Get("/v2/bot/channel/webhook/endpoint", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var webhookSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set webhook endpoint URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		endpoint, _ := cmd.Flags().GetString("url")
		if endpoint == "" {
			return fmt.Errorf("--url is required")
		}

		c := client.New(token)
		data, err := c.Put("/v2/bot/channel/webhook/endpoint", map[string]any{
			"endpoint": endpoint,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("webhook endpoint updated")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	webhookSetCmd.Flags().String("url", "", "Webhook endpoint URL")
	webhookSetCmd.MarkFlagRequired("url")

	webhookCmd.AddCommand(webhookGetCmd)
	webhookCmd.AddCommand(webhookSetCmd)
}
