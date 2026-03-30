package cmd

import (
	"fmt"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage group and room chats",
}

var groupMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Get group member user IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		groupID, _ := cmd.Flags().GetString("group-id")
		if groupID == "" {
			return fmt.Errorf("--group-id is required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/group/"+groupID+"/members/ids", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var groupMemberProfileCmd = &cobra.Command{
	Use:   "member-profile",
	Short: "Get a group member's profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		groupID, _ := cmd.Flags().GetString("group-id")
		userID, _ := cmd.Flags().GetString("user-id")
		if groupID == "" || userID == "" {
			return fmt.Errorf("--group-id and --user-id are required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/group/"+groupID+"/member/"+userID, nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var groupLeaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Leave a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		groupID, _ := cmd.Flags().GetString("group-id")
		if groupID == "" {
			return fmt.Errorf("--group-id is required")
		}

		c := client.New(token)
		_, err := c.Post("/v2/bot/group/"+groupID+"/leave", nil)
		if err != nil {
			return err
		}
		fmt.Println("left group")
		return nil
	},
}

func init() {
	groupMembersCmd.Flags().String("group-id", "", "Group ID")
	groupMemberProfileCmd.Flags().String("group-id", "", "Group ID")
	groupMemberProfileCmd.Flags().String("user-id", "", "User ID")
	groupLeaveCmd.Flags().String("group-id", "", "Group ID")

	groupCmd.AddCommand(groupMembersCmd)
	groupCmd.AddCommand(groupMemberProfileCmd)
	groupCmd.AddCommand(groupLeaveCmd)
}
