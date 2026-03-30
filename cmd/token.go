package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage channel access tokens",
}

var tokenIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue a short-lived channel access token",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")

		if clientID == "" {
			clientID = config.Get(config.KeyClientID)
		}
		if clientSecret == "" {
			clientSecret = config.Get(config.KeyClientSecret)
		}
		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("--client-id and --client-secret are required (or set via: line config set)")
		}

		c := client.New("")
		data, err := c.PostForm("https://api.line.me/v2/oauth/accessToken", map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": clientSecret,
		})
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var tokenVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a channel access token",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = config.RequireToken()
		}

		c := client.New("")
		data, err := c.PostForm("https://api.line.me/v2/oauth/verify", map[string]string{
			"access_token": token,
		})
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke a channel access token",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = config.RequireToken()
		}

		c := client.New("")
		data, err := c.PostForm("https://api.line.me/v2/oauth/revoke", map[string]string{
			"access_token": token,
		})
		if err != nil {
			return err
		}
		if len(data) == 0 {
			fmt.Println("token revoked")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	tokenIssueCmd.Flags().String("client-id", "", "Channel client ID")
	tokenIssueCmd.Flags().String("client-secret", "", "Channel client secret")

	tokenVerifyCmd.Flags().String("token", "", "Access token to verify (defaults to configured token)")
	tokenRevokeCmd.Flags().String("token", "", "Access token to revoke (defaults to configured token)")

	tokenCmd.AddCommand(tokenIssueCmd)
	tokenCmd.AddCommand(tokenVerifyCmd)
	tokenCmd.AddCommand(tokenRevokeCmd)
}
