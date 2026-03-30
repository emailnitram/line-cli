package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if token == "" && clientID == "" && clientSecret == "" && apiKey == "" {
			return fmt.Errorf("provide at least one flag: --token, --client-id, --client-secret, --api-key")
		}

		if token != "" {
			if err := config.Set(config.KeyAccessToken, token); err != nil {
				return err
			}
			fmt.Println("access token saved")
		}
		if clientID != "" {
			if err := config.Set(config.KeyClientID, clientID); err != nil {
				return err
			}
			fmt.Println("client id saved")
		}
		if clientSecret != "" {
			if err := config.Set(config.KeyClientSecret, clientSecret); err != nil {
				return err
			}
			fmt.Println("client secret saved")
		}
		if apiKey != "" {
			if err := config.Set(config.KeyAPIKey, apiKey); err != nil {
				return err
			}
			fmt.Println("api key saved")
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		token := config.Get(config.KeyAccessToken)
		clientID := config.Get(config.KeyClientID)
		clientSecret := config.Get(config.KeyClientSecret)
		apiKey := config.Get(config.KeyAPIKey)

		mask := func(s string) string {
			if len(s) <= 8 {
				return "****"
			}
			return s[:4] + "..." + s[len(s)-4:]
		}

		if token != "" {
			fmt.Printf("access_token:  %s\n", mask(token))
		} else {
			fmt.Println("access_token:  (not set)")
		}
		if clientID != "" {
			fmt.Printf("client_id:     %s\n", clientID)
		} else {
			fmt.Println("client_id:     (not set)")
		}
		if clientSecret != "" {
			fmt.Printf("client_secret: %s\n", mask(clientSecret))
		} else {
			fmt.Println("client_secret: (not set)")
		}
		if apiKey != "" {
			fmt.Printf("api_key:       %s\n", mask(apiKey))
		} else {
			fmt.Println("api_key:       (not set)")
		}
	},
}

func init() {
	configSetCmd.Flags().String("token", "", "Channel access token")
	configSetCmd.Flags().String("client-id", "", "Channel client ID")
	configSetCmd.Flags().String("client-secret", "", "Channel client secret")
	configSetCmd.Flags().String("api-key", "", "LINE Shopping API key")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}
