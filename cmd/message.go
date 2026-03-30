package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Send messages",
}

var messagePushCmd = &cobra.Command{
	Use:   "push",
	Short: "Send a push message to a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		to, _ := cmd.Flags().GetString("to")
		text, _ := cmd.Flags().GetString("text")
		file, _ := cmd.Flags().GetString("file")

		if to == "" {
			return fmt.Errorf("--to is required")
		}

		var messages []map[string]any
		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(data, &messages); err != nil {
				return err
			}
		} else if text != "" {
			messages = []map[string]any{{"type": "text", "text": text}}
		} else {
			return fmt.Errorf("provide --text or --file")
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/message/push", map[string]any{
			"to":       to,
			"messages": messages,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("message sent")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var messageBroadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "Send a broadcast message to all followers",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		text, _ := cmd.Flags().GetString("text")
		file, _ := cmd.Flags().GetString("file")

		var messages []map[string]any
		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(data, &messages); err != nil {
				return err
			}
		} else if text != "" {
			messages = []map[string]any{{"type": "text", "text": text}}
		} else {
			return fmt.Errorf("provide --text or --file")
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/message/broadcast", map[string]any{
			"messages": messages,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("broadcast sent")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var messageReplyCmd = &cobra.Command{
	Use:   "reply",
	Short: "Reply to a message using a reply token",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		replyToken, _ := cmd.Flags().GetString("reply-token")
		text, _ := cmd.Flags().GetString("text")
		file, _ := cmd.Flags().GetString("file")

		if replyToken == "" {
			return fmt.Errorf("--reply-token is required")
		}

		var messages []map[string]any
		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(data, &messages); err != nil {
				return err
			}
		} else if text != "" {
			messages = []map[string]any{{"type": "text", "text": text}}
		} else {
			return fmt.Errorf("provide --text or --file")
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/message/reply", map[string]any{
			"replyToken": replyToken,
			"messages":   messages,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("reply sent")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var messageQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Get monthly message quota and usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		data, err := c.Get("/v2/bot/message/quota", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var messageLoadingCmd = &cobra.Command{
	Use:   "loading",
	Short: "Display a loading animation in a chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		chatID, _ := cmd.Flags().GetString("chat-id")
		seconds, _ := cmd.Flags().GetInt("seconds")

		if chatID == "" {
			return fmt.Errorf("--chat-id is required")
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/chat/loading/start", map[string]any{
			"chatId":         chatID,
			"loadingSeconds": seconds,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("loading animation started")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	messagePushCmd.Flags().String("to", "", "User ID to send to")
	messagePushCmd.Flags().String("text", "", "Text message content")
	messagePushCmd.Flags().String("file", "", "Path to JSON file containing messages array")

	messageBroadcastCmd.Flags().String("text", "", "Text message content")
	messageBroadcastCmd.Flags().String("file", "", "Path to JSON file containing messages array")

	messageReplyCmd.Flags().String("reply-token", "", "Reply token from webhook event")
	messageReplyCmd.Flags().String("text", "", "Text message content")
	messageReplyCmd.Flags().String("file", "", "Path to JSON file containing messages array")

	messageLoadingCmd.Flags().String("chat-id", "", "User or group ID")
	messageLoadingCmd.Flags().Int("seconds", 5, "Loading animation duration (5-60)")

	messageCmd.AddCommand(messagePushCmd)
	messageCmd.AddCommand(messageBroadcastCmd)
	messageCmd.AddCommand(messageReplyCmd)
	messageCmd.AddCommand(messageQuotaCmd)
	messageCmd.AddCommand(messageLoadingCmd)
}
