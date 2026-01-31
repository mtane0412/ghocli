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

// ImagesUploadCmd は画像をアップロードするコマンドです
type ImagesUploadCmd struct {
	File    string `arg:"" help:"Path to image file" type:"existingfile"`
	Purpose string `help:"Image purpose (image, profile_image, icon)" short:"p" default:"image"`
	Ref     string `help:"Reference ID for the image" short:"r"`
}

// Run はimagesコマンドのuploadサブコマンドを実行します
func (c *ImagesUploadCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ファイルを開く
	file, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("ファイルのオープンに失敗: %w", err)
	}
	defer file.Close()

	// ファイル名を取得
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ファイル情報の取得に失敗: %w", err)
	}

	// 画像をアップロード
	image, err := client.UploadImage(file, fileInfo.Name(), ghostapi.ImageUploadOptions{
		Purpose: c.Purpose,
		Ref:     c.Ref,
	})
	if err != nil {
		return fmt.Errorf("画像のアップロードに失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(image)
	}

	// 成功メッセージと画像URLを表示
	formatter.PrintMessage(fmt.Sprintf("画像をアップロードしました: %s", image.URL))

	return nil
}
