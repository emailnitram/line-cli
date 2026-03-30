package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var couponCmd = &cobra.Command{
	Use:   "coupon",
	Short: "Manage coupons",
}

var couponCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a coupon from a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			return fmt.Errorf("--file is required")
		}

		raw, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var body any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/coupon", body)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var couponListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of coupons",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		limit, _ := cmd.Flags().GetString("limit")
		status, _ := cmd.Flags().GetString("status")

		params := map[string]string{}
		if limit != "" {
			params["limit"] = limit
		}
		if status != "" {
			params["status"] = status
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/coupon", params)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var couponGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a coupon",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/coupon/"+id, nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var couponDiscontinueCmd = &cobra.Command{
	Use:   "discontinue",
	Short: "Discontinue a coupon",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		_, err := c.Put("/v2/bot/coupon/"+id+"/close", nil)
		if err != nil {
			return err
		}
		fmt.Println("coupon discontinued")
		return nil
	},
}

var couponSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a coupon as a broadcast message",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		to, _ := cmd.Flags().GetString("to")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		messages := []map[string]any{{"type": "coupon", "couponId": id}}

		c := client.New(token)
		var (
			data []byte
			err  error
		)
		if to != "" {
			data, err = c.Post("/v2/bot/message/push", map[string]any{
				"to":       to,
				"messages": messages,
			})
		} else {
			data, err = c.Post("/v2/bot/message/broadcast", map[string]any{
				"messages": messages,
			})
		}
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("coupon sent")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	couponCreateCmd.Flags().String("file", "", "Path to coupon JSON file")
	couponListCmd.Flags().String("limit", "100", "Max number of coupons to return")
	couponListCmd.Flags().String("status", "", "Filter by status (DRAFT, RUNNING, CLOSED)")
	couponGetCmd.Flags().String("id", "", "Coupon ID")
	couponDiscontinueCmd.Flags().String("id", "", "Coupon ID")
	couponSendCmd.Flags().String("id", "", "Coupon ID")
	couponSendCmd.Flags().String("to", "", "User ID (omit to broadcast)")

	couponCmd.AddCommand(couponCreateCmd)
	couponCmd.AddCommand(couponListCmd)
	couponCmd.AddCommand(couponGetCmd)
	couponCmd.AddCommand(couponDiscontinueCmd)
	couponCmd.AddCommand(couponSendCmd)
}
