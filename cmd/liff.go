package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var liffCmd = &cobra.Command{
	Use:   "liff",
	Short: "Manage LIFF apps",
}

var liffListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all LIFF apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		data, err := c.Get("/liff/v1/apps", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var liffCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a LIFF app",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		url, _ := cmd.Flags().GetString("url")
		viewType, _ := cmd.Flags().GetString("type")
		file, _ := cmd.Flags().GetString("file")

		var body any

		if file != "" {
			raw, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(raw, &body); err != nil {
				return err
			}
		} else {
			if url == "" || viewType == "" {
				return fmt.Errorf("--url and --type are required (or use --file)")
			}
			validTypes := map[string]bool{"compact": true, "tall": true, "full": true}
			if !validTypes[viewType] {
				return fmt.Errorf("--type must be one of: compact, tall, full")
			}
			body = map[string]any{
				"view": map[string]any{
					"type": viewType,
					"url":  url,
				},
			}
		}

		c := client.New(token)
		data, err := c.Post("/liff/v1/apps", body)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var liffUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a LIFF app",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		url, _ := cmd.Flags().GetString("url")
		viewType, _ := cmd.Flags().GetString("type")
		file, _ := cmd.Flags().GetString("file")

		if id == "" {
			return fmt.Errorf("--id is required")
		}

		var body any

		if file != "" {
			raw, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(raw, &body); err != nil {
				return err
			}
		} else {
			if url == "" && viewType == "" {
				return fmt.Errorf("provide at least --url or --type (or use --file)")
			}
			view := map[string]any{}
			if url != "" {
				view["url"] = url
			}
			if viewType != "" {
				validTypes := map[string]bool{"compact": true, "tall": true, "full": true}
				if !validTypes[viewType] {
					return fmt.Errorf("--type must be one of: compact, tall, full")
				}
				view["type"] = viewType
			}
			body = map[string]any{"view": view}
		}

		c := client.New(token)
		_, err := c.Put("/liff/v1/apps/"+id, body)
		if err != nil {
			return err
		}
		fmt.Printf("LIFF app %s updated\n", id)
		return nil
	},
}

var liffDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a LIFF app",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		_, err := c.Delete("/liff/v1/apps/" + id)
		if err != nil {
			return err
		}
		fmt.Printf("LIFF app %s deleted\n", id)
		return nil
	},
}

func init() {
	liffCreateCmd.Flags().String("url", "", "Endpoint URL of the LIFF app")
	liffCreateCmd.Flags().String("type", "", "View type: compact, tall, full")
	liffCreateCmd.Flags().String("file", "", "Path to JSON file with full LIFF app config")

	liffUpdateCmd.Flags().String("id", "", "LIFF app ID")
	liffUpdateCmd.Flags().String("url", "", "New endpoint URL")
	liffUpdateCmd.Flags().String("type", "", "New view type: compact, tall, full")
	liffUpdateCmd.Flags().String("file", "", "Path to JSON file with full LIFF app config")

	liffDeleteCmd.Flags().String("id", "", "LIFF app ID")

	liffCmd.AddCommand(liffListCmd)
	liffCmd.AddCommand(liffCreateCmd)
	liffCmd.AddCommand(liffUpdateCmd)
	liffCmd.AddCommand(liffDeleteCmd)
}
