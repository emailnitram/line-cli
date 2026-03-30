package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/emailnitram/line-cli/internal/client"
	"github.com/emailnitram/line-cli/internal/config"
	"github.com/spf13/cobra"
)

// validSizes are the LINE-accepted rich menu dimensions (width x height).
var validSizes = [][2]int{
	{2500, 1686},
	{2500, 843},
	{1200, 810},
	{1200, 405},
	{800, 540},
	{800, 270},
}

const maxImageBytes = 1 * 1024 * 1024 // 1 MB

// prepareRichMenuImage resizes the image to the target dimensions, compresses
// it as JPEG, and reduces quality until it fits under 1 MB.
func prepareRichMenuImage(file string, width, height int) ([]byte, error) {
	src, err := imaging.Open(file, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}

	resized := imaging.Resize(src, width, height, imaging.Lanczos)

	for quality := 85; quality >= 40; quality -= 10 {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}
		if buf.Len() <= maxImageBytes {
			fmt.Printf("image prepared: %dx%d, %dKB (quality %d)\n", width, height, buf.Len()/1024, quality)
			return buf.Bytes(), nil
		}
	}
	return nil, fmt.Errorf("image could not be compressed under 1MB at minimum quality")
}

var richMenuCmd = &cobra.Command{
	Use:   "richmenu",
	Short: "Manage rich menus",
}

var richMenuListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rich menus",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		data, err := c.Get("/v2/bot/richmenu/list", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var richMenuCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a rich menu from a JSON file",
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
		data, err := c.Post("/v2/bot/richmenu", body)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var richMenuDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a rich menu",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		_, err := c.Delete("/v2/bot/richmenu/" + id)
		if err != nil {
			return err
		}
		fmt.Println("rich menu deleted")
		return nil
	},
}

var richMenuSetDefaultCmd = &cobra.Command{
	Use:   "set-default",
	Short: "Set the default rich menu",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}

		c := client.New(token)
		_, err := c.Post("/v2/bot/user/all/richmenu/"+id, nil)
		if err != nil {
			return err
		}
		fmt.Println("default rich menu set")
		return nil
	},
}

var richMenuCancelDefaultCmd = &cobra.Command{
	Use:   "cancel-default",
	Short: "Cancel the default rich menu",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		c := client.New(token)
		_, err := c.Delete("/v2/bot/user/all/richmenu")
		if err != nil {
			return err
		}
		fmt.Println("default rich menu cancelled")
		return nil
	},
}

var richMenuLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link a rich menu to a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		userID, _ := cmd.Flags().GetString("user-id")
		if id == "" || userID == "" {
			return fmt.Errorf("--id and --user-id are required")
		}

		c := client.New(token)
		_, err := c.Post("/v2/bot/user/"+userID+"/richmenu/"+id, nil)
		if err != nil {
			return err
		}
		fmt.Printf("rich menu %s linked to user %s\n", id, userID)
		return nil
	},
}

var richMenuGetUserCmd = &cobra.Command{
	Use:   "get-user",
	Short: "Get the rich menu linked to a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		userID, _ := cmd.Flags().GetString("user-id")
		if userID == "" {
			return fmt.Errorf("--user-id is required")
		}

		c := client.New(token)
		data, err := c.Get("/v2/bot/user/"+userID+"/richmenu", nil)
		if err != nil {
			return err
		}
		client.PrintJSON(data)
		return nil
	},
}

var richMenuUploadImageCmd = &cobra.Command{
	Use:   "upload-image",
	Short: "Upload an image for a rich menu (auto-resizes and compresses)",
	Long: `Upload an image for a rich menu.

The image is automatically resized to the target dimensions and compressed to
fit LINE's 1MB limit. Valid sizes (WxH):
  2500x1686  2500x843 (default)
  1200x810   1200x405
  800x540    800x270`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		id, _ := cmd.Flags().GetString("id")
		file, _ := cmd.Flags().GetString("file")
		size, _ := cmd.Flags().GetString("size")
		if id == "" || file == "" {
			return fmt.Errorf("--id and --file are required")
		}

		// Parse --size (WxH)
		width, height := 2500, 843
		if size != "" {
			if _, err := fmt.Sscanf(size, "%dx%d", &width, &height); err != nil {
				return fmt.Errorf("--size must be in WxH format, e.g. 2500x843")
			}
			valid := false
			for _, s := range validSizes {
				if s[0] == width && s[1] == height {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid size %dx%d — valid sizes: 2500x1686, 2500x843, 1200x810, 1200x405, 800x540, 800x270", width, height)
			}
		}

		imgData, err := prepareRichMenuImage(file, width, height)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost,
			client.BaseDataURL+"/v2/bot/richmenu/"+id+"/content",
			strings.NewReader(string(imgData)))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "image/jpeg")
		req.ContentLength = int64(len(imgData))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		fmt.Printf("image uploaded to rich menu %s\n", id)
		return nil
	},
}

var richMenuCreateAliasCmd = &cobra.Command{
	Use:   "create-alias",
	Short: "Create a rich menu alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.RequireToken()
		aliasID, _ := cmd.Flags().GetString("alias-id")
		menuID, _ := cmd.Flags().GetString("menu-id")
		if aliasID == "" || menuID == "" {
			return fmt.Errorf("--alias-id and --menu-id are required")
		}

		c := client.New(token)
		data, err := c.Post("/v2/bot/richmenu/alias", map[string]any{
			"richMenuAliasId": aliasID,
			"richMenuId":      menuID,
		})
		if err != nil {
			return err
		}
		if len(data) <= 2 {
			fmt.Println("alias created")
			return nil
		}
		client.PrintJSON(data)
		return nil
	},
}

func init() {
	richMenuCreateCmd.Flags().String("file", "", "Path to rich menu JSON file")
	richMenuDeleteCmd.Flags().String("id", "", "Rich menu ID")
	richMenuSetDefaultCmd.Flags().String("id", "", "Rich menu ID")
	richMenuLinkCmd.Flags().String("id", "", "Rich menu ID")
	richMenuLinkCmd.Flags().String("user-id", "", "User ID")
	richMenuGetUserCmd.Flags().String("user-id", "", "User ID")
	richMenuCreateAliasCmd.Flags().String("alias-id", "", "Rich menu alias ID")
	richMenuCreateAliasCmd.Flags().String("menu-id", "", "Rich menu ID")
	richMenuUploadImageCmd.Flags().String("id", "", "Rich menu ID")
	richMenuUploadImageCmd.Flags().String("file", "", "Path to image file (JPEG or PNG, any size)")
	richMenuUploadImageCmd.Flags().String("size", "2500x843", "Target dimensions WxH (2500x1686, 2500x843, 1200x810, 1200x405, 800x540, 800x270)")

	richMenuCmd.AddCommand(richMenuListCmd)
	richMenuCmd.AddCommand(richMenuCreateCmd)
	richMenuCmd.AddCommand(richMenuDeleteCmd)
	richMenuCmd.AddCommand(richMenuSetDefaultCmd)
	richMenuCmd.AddCommand(richMenuCancelDefaultCmd)
	richMenuCmd.AddCommand(richMenuLinkCmd)
	richMenuCmd.AddCommand(richMenuGetUserCmd)
	richMenuCmd.AddCommand(richMenuCreateAliasCmd)
	richMenuCmd.AddCommand(richMenuUploadImageCmd)
}
