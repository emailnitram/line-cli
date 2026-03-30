package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var audienceCmd = &cobra.Command{
	Use:   "audience",
	Short: "Manage audiences",
}

var audienceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an audience",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		_, err := c.Delete("/v2/bot/audienceGroup/" + id)
		if err != nil {
			return err
		}
		fmt.Println("audience deleted")
		return nil
	},
}

var audienceRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename an audience",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		if id == "" || name == "" {
			return fmt.Errorf("--id and --name are required")
		}

		c := client.New(token)
		data, err := c.Put("/v2/bot/audienceGroup/"+id+"/updateDescription", map[string]any{
			"description": name,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("audience renamed")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var audienceGetSharedCmd = &cobra.Command{
	Use:   "get-shared",
	Short: "Get shared audience data from Business Manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/audienceGroup/shared/"+id, nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	audienceDeleteCmd.Flags().String("id", "", "Audience ID")
	audienceRenameCmd.Flags().String("id", "", "Audience ID")
	audienceRenameCmd.Flags().String("name", "", "New audience name")
	audienceGetSharedCmd.Flags().String("id", "", "Audience group ID")

	audienceCmd.AddCommand(audienceDeleteCmd)
	audienceCmd.AddCommand(audienceRenameCmd)
	audienceCmd.AddCommand(audienceGetSharedCmd)
}
