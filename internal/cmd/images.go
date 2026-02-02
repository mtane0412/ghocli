/**
 * images.go
 * 画像管理コマンド
 *
 * Ghost画像のアップロード機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// ImagesCmd は画像管理コマンドです
type ImagesCmd struct {
	Upload ImagesUploadCmd `cmd:"" help:"Upload an image"`
}

// ImagesUploadCmd is the command to upload 画像
type ImagesUploadCmd struct {
	File    string `arg:"" help:"Path to image file" type:"existingfile"`
	Purpose string `help:"Image purpose (image, profile_image, icon)" short:"p" default:"image"`
	Ref     string `help:"Reference ID for the image" short:"r"`
}

// Run executes the upload subcommand of the images command
func (c *ImagesUploadCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Open file
	file, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get filename
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file information: %w", err)
	}

	// Upload image
	image, err := client.UploadImage(file, fileInfo.Name(), ghostapi.ImageUploadOptions{
		Purpose: c.Purpose,
		Ref:     c.Ref,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(image)
	}

	// 成功メッセージと画像URLを表示
	formatter.PrintMessage(fmt.Sprintf("uploaded image: %s", image.URL))

	return nil
}
