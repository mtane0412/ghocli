/**
 * themes.go
 * テーマ管理コマンド
 *
 * Ghostテーマの一覧表示、アップロード、有効化機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mtane0412/gho/internal/outfmt"
)

// ThemesCmd はテーマ管理コマンドです
type ThemesCmd struct {
	List     ThemesListCmd     `cmd:"" help:"List themes"`
	Upload   ThemesUploadCmd   `cmd:"" help:"Upload a theme"`
	Activate ThemesActivateCmd `cmd:"" help:"Activate a theme"`
	Delete   ThemesDeleteCmd   `cmd:"" help:"Delete a theme"`

	// Phase 3: 複合操作
	Install ThemesInstallCmd `cmd:"" help:"Upload and activate a theme"`
}

// ThemesListCmd はテーマ一覧を取得するコマンドです
type ThemesListCmd struct{}

// Run はthemesコマンドのlistサブコマンドを実行します
func (c *ThemesListCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// テーマ一覧を取得
	response, err := client.ListThemes()
	if err != nil {
		return fmt.Errorf("テーマ一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Themes)
	}

	// テーブル形式で出力
	headers := []string{"Name", "Active", "Version", "Description"}
	rows := make([][]string, len(response.Themes))
	for i, theme := range response.Themes {
		active := ""
		if theme.Active {
			active = "✓"
		}

		version := ""
		description := ""
		if theme.Package != nil {
			version = theme.Package.Version
			description = theme.Package.Description
		}

		rows[i] = []string{
			theme.Name,
			active,
			version,
			description,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// ThemesUploadCmd はテーマをアップロードするコマンドです
type ThemesUploadCmd struct {
	File string `arg:"" help:"Theme zip file path" type:"existingfile"`
}

// Run はthemesコマンドのuploadサブコマンドを実行します
func (c *ThemesUploadCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// テーマをアップロード
	theme, err := client.UploadTheme(file, fileInfo.Name())
	if err != nil {
		return fmt.Errorf("テーマのアップロードに失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("テーマをアップロードしました: %s", theme.Name))
		if theme.Package != nil && theme.Package.Version != "" {
			formatter.PrintMessage(fmt.Sprintf("バージョン: %s", theme.Package.Version))
		}
	}

	// JSON形式の場合はテーマ情報も出力
	if root.JSON {
		return formatter.Print(theme)
	}

	return nil
}

// ThemesActivateCmd はテーマを有効化するコマンドです
type ThemesActivateCmd struct {
	Name string `arg:"" help:"Theme name"`
}

// Run はthemesコマンドのactivateサブコマンドを実行します
func (c *ThemesActivateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// テーマを有効化
	theme, err := client.ActivateTheme(c.Name)
	if err != nil {
		return fmt.Errorf("テーマの有効化に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("テーマを有効化しました: %s", theme.Name))
	}

	// JSON形式の場合はテーマ情報も出力
	if root.JSON {
		return formatter.Print(theme)
	}

	return nil
}

// ========================================
// Phase 3: 複合操作
// ========================================

// ThemesInstallCmd はテーマをアップロードして有効化するコマンドです
type ThemesInstallCmd struct {
	File string `arg:"" help:"Path to theme zip file" type:"existingfile"`
}

// Run はthemesコマンドのinstallサブコマンドを実行します
func (c *ThemesInstallCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// ファイルを開く
	file, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("ファイルを開くことに失敗: %w", err)
	}
	defer file.Close()

	// ファイル名を取得
	filename := filepath.Base(c.File)

	// テーマをアップロード
	formatter.PrintMessage(fmt.Sprintf("テーマをアップロード中: %s", c.File))
	uploadedTheme, err := client.UploadTheme(file, filename)
	if err != nil {
		return fmt.Errorf("テーマのアップロードに失敗: %w", err)
	}

	formatter.PrintMessage(fmt.Sprintf("テーマをアップロードしました: %s", uploadedTheme.Name))

	// テーマを有効化
	formatter.PrintMessage(fmt.Sprintf("テーマを有効化中: %s", uploadedTheme.Name))
	activatedTheme, err := client.ActivateTheme(uploadedTheme.Name)
	if err != nil {
		return fmt.Errorf("テーマの有効化に失敗: %w", err)
	}

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("テーマをインストールして有効化しました: %s", activatedTheme.Name))
	}

	// JSON形式の場合はテーマ情報も出力
	if root.JSON {
		return formatter.Print(activatedTheme)
	}

	return nil
}

// ========================================
// テーマ削除
// ========================================

// ThemesDeleteCmd はテーマを削除するコマンドです
type ThemesDeleteCmd struct {
	Name string `arg:"" help:"Theme name"`
}

// Run はthemesコマンドのdeleteサブコマンドを実行します
func (c *ThemesDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// テーマ一覧を取得してアクティブかチェック
	themes, err := client.ListThemes()
	if err != nil {
		return fmt.Errorf("テーマ一覧の取得に失敗: %w", err)
	}

	// アクティブなテーマの削除を防止
	for _, theme := range themes.Themes {
		if theme.Name == c.Name && theme.Active {
			return fmt.Errorf("アクティブなテーマは削除できません: %s", c.Name)
		}
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete theme '%s'", c.Name)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// テーマを削除
	if err := client.DeleteTheme(c.Name); err != nil {
		return fmt.Errorf("テーマの削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("テーマを削除しました: %s", c.Name))

	return nil
}
