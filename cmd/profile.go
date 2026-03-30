package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Get user profile information",
}

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a user's profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		userID, _ := cmd.Flags().GetString("user-id")
		if userID == "" {
			return fmt.Errorf("--user-id is required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/profile/"+userID, nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var profileDemographicsCmd = &cobra.Command{
	Use:   "demographics",
	Short: "Get friend demographics",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		data, err := c.Get("/v2/bot/insight/demographic", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	profileGetCmd.Flags().String("user-id", "", "LINE user ID")
	profileGetCmd.MarkFlagRequired("user-id")

	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileDemographicsCmd)
}
