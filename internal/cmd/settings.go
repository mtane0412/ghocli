/**
 * settings.go
 * 設定管理コマンド
 *
 * Ghostサイト設定の表示、更新機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// SettingsCmd は設定管理コマンドです
type SettingsCmd struct {
	List SettingsListCmd `cmd:"" help:"List all settings"`
	Get  SettingsGetCmd  `cmd:"" help:"Get a specific setting"`
	Set  SettingsSetCmd  `cmd:"" help:"Set a setting value"`
}

// SettingsListCmd は設定一覧を取得するコマンドです
type SettingsListCmd struct{}

// Run はsettingsコマンドのlistサブコマンドを実行します
func (c *SettingsListCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 設定一覧を取得
	response, err := client.GetSettings()
	if err != nil {
		return fmt.Errorf("設定一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Settings)
	}

	// テーブル形式で出力
	headers := []string{"Key", "Value"}
	rows := make([][]string, len(response.Settings))
	for i, setting := range response.Settings {
		value := fmt.Sprintf("%v", setting.Value)
		// 長い値は切り詰める
		if len(value) > 80 {
			value = value[:77] + "..."
		}
		rows[i] = []string{
			setting.Key,
			value,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// SettingsGetCmd は特定の設定値を取得するコマンドです
type SettingsGetCmd struct {
	Key string `arg:"" help:"Setting key"`
}

// Run はsettingsコマンドのgetサブコマンドを実行します
func (c *SettingsGetCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 設定一覧を取得
	response, err := client.GetSettings()
	if err != nil {
		return fmt.Errorf("設定一覧の取得に失敗: %w", err)
	}

	// 指定されたキーの設定を検索
	var foundSetting *ghostapi.Setting
	for _, setting := range response.Settings {
		if setting.Key == c.Key {
			foundSetting = &setting
			break
		}
	}

	if foundSetting == nil {
		return fmt.Errorf("設定が見つかりません: %s", c.Key)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(foundSetting)
	}

	// 値を出力
	formatter.PrintMessage(fmt.Sprintf("%v", foundSetting.Value))

	return nil
}

// SettingsSetCmd は設定値を更新するコマンドです
type SettingsSetCmd struct {
	Key   string `arg:"" help:"Setting key"`
	Value string `arg:"" help:"Setting value"`
}

// Run はsettingsコマンドのsetサブコマンドを実行します
func (c *SettingsSetCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 設定を更新
	updates := []ghostapi.SettingUpdate{
		{
			Key:   c.Key,
			Value: c.Value,
		},
	}

	response, err := client.UpdateSettings(updates)
	if err != nil {
		return fmt.Errorf("設定の更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("設定を更新しました: %s = %s", c.Key, c.Value))
	}

	// JSON形式の場合は設定情報も出力
	if root.JSON {
		return formatter.Print(response.Settings)
	}

	return nil
}
