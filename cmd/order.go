package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Manage LINE Shopping orders",
}

func shopRequest(method, path string, body any) ([]byte, error) {
	apiKey := config.Get(config.KeyAPIKey)
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: no api key set. Run: line config set --api-key <key>")
		os.Exit(1)
	}
	company := os.Getenv("LINE_COMPANY")
	if company == "" {
		company = "line-cli"
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, client.BaseOAURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", company)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

var orderGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get order details",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderNo, _ := cmd.Flags().GetString("order-no")
		if orderNo == "" {
			return fmt.Errorf("--order-no is required")
		}

		data, err := shopRequest(http.MethodGet, "/myshop/v1/orders/"+orderNo, nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var orderMarkPaidCmd = &cobra.Command{
	Use:   "mark-paid",
	Short: "Mark a COD order as paid",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderNo, _ := cmd.Flags().GetString("order-no")
		if orderNo == "" {
			return fmt.Errorf("--order-no is required")
		}

		data, err := shopRequest(http.MethodPost, "/myshop/v1/orders/"+orderNo+"/paid", nil)
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("order marked as paid")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var orderMarkShippedCmd = &cobra.Command{
	Use:   "mark-shipped",
	Short: "Mark an order as shipped",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderNo, _ := cmd.Flags().GetString("order-no")
		if orderNo == "" {
			return fmt.Errorf("--order-no is required")
		}

		data, err := shopRequest(http.MethodPost, "/myshop/v1/orders/"+orderNo+"/ship", nil)
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("order marked as shipped")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

var orderTrackingCmd = &cobra.Command{
	Use:   "tracking",
	Short: "Update tracking number for an order",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderNo, _ := cmd.Flags().GetString("order-no")
		tracking, _ := cmd.Flags().GetString("tracking-no")
		if orderNo == "" || tracking == "" {
			return fmt.Errorf("--order-no and --tracking-no are required")
		}

		data, err := shopRequest(http.MethodPut, "/myshop/v1/orders/"+orderNo+"/tracking", map[string]any{
			"trackingNumber": tracking,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("tracking number updated")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	orderGetCmd.Flags().String("order-no", "", "Order number")
	orderMarkPaidCmd.Flags().String("order-no", "", "Order number")
	orderMarkShippedCmd.Flags().String("order-no", "", "Order number")
	orderTrackingCmd.Flags().String("order-no", "", "Order number")
	orderTrackingCmd.Flags().String("tracking-no", "", "Tracking number")

	orderCmd.AddCommand(orderGetCmd)
	orderCmd.AddCommand(orderMarkPaidCmd)
	orderCmd.AddCommand(orderMarkShippedCmd)
	orderCmd.AddCommand(orderTrackingCmd)
}
