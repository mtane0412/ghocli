/**
 * webhooks.go
 * Webhook管理コマンド
 *
 * GhostのWebhookの作成、更新、削除機能を提供します。
 * 注意: Ghost APIはWebhookのList/Getをサポートしていません。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// WebhooksCmd はWebhook管理コマンドです
type WebhooksCmd struct {
	Create WebhooksCreateCmd `cmd:"" help:"Create a webhook"`
	Update WebhooksUpdateCmd `cmd:"" help:"Update a webhook"`
	Delete WebhooksDeleteCmd `cmd:"" help:"Delete a webhook"`
}

// WebhooksCreateCmd はWebhookを作成するコマンドです
type WebhooksCreateCmd struct {
	Event     string `help:"Webhook event (e.g., post.published, member.added)" short:"e" required:""`
	TargetURL string `help:"Target URL for webhook" short:"t" required:""`
	Name      string `help:"Webhook name" short:"n"`
}

// Run はwebhooksコマンドのcreateサブコマンドを実行します
func (c *WebhooksCreateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 新規Webhookを作成
	webhook := &ghostapi.Webhook{
		Event:     c.Event,
		TargetURL: c.TargetURL,
		Name:      c.Name,
	}

	created, err := client.CreateWebhook(webhook)
	if err != nil {
		return fmt.Errorf("webhookの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("Webhookを作成しました (ID: %s)", created.ID))
		formatter.PrintMessage(fmt.Sprintf("イベント: %s", created.Event))
		formatter.PrintMessage(fmt.Sprintf("URL: %s", created.TargetURL))
		if created.Secret != "" {
			formatter.PrintMessage(fmt.Sprintf("シークレット: %s", created.Secret))
		}
	}

	// JSON形式の場合はWebhook情報も出力
	if root.JSON {
		return formatter.Print(created)
	}

	return nil
}

// WebhooksUpdateCmd はWebhookを更新するコマンドです
type WebhooksUpdateCmd struct {
	ID        string `arg:"" help:"Webhook ID"`
	Event     string `help:"Webhook event" short:"e"`
	TargetURL string `help:"Target URL for webhook" short:"t"`
	Name      string `help:"Webhook name" short:"n"`
}

// Run はwebhooksコマンドのupdateサブコマンドを実行します
func (c *WebhooksUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 更新内容を作成（指定されたフィールドのみ）
	webhook := &ghostapi.Webhook{}

	if c.Event != "" {
		webhook.Event = c.Event
	}
	if c.TargetURL != "" {
		webhook.TargetURL = c.TargetURL
	}
	if c.Name != "" {
		webhook.Name = c.Name
	}

	// Webhookを更新
	updated, err := client.UpdateWebhook(c.ID, webhook)
	if err != nil {
		return fmt.Errorf("webhookの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("Webhookを更新しました (ID: %s)", updated.ID))
		formatter.PrintMessage(fmt.Sprintf("イベント: %s", updated.Event))
		formatter.PrintMessage(fmt.Sprintf("URL: %s", updated.TargetURL))
	}

	// JSON形式の場合はWebhook情報も出力
	if root.JSON {
		return formatter.Print(updated)
	}

	return nil
}

// WebhooksDeleteCmd はWebhookを削除するコマンドです
type WebhooksDeleteCmd struct {
	ID string `arg:"" help:"Webhook ID"`
}

// Run はwebhooksコマンドのdeleteサブコマンドを実行します
func (c *WebhooksDeleteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 確認なしで削除する場合を除き、確認を求める
	if !root.Force {
		fmt.Printf("本当にWebhook (ID: %s)を削除しますか? [y/N]: ", c.ID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("削除がキャンセルされました")
		}
	}

	// Webhookを削除
	if err := client.DeleteWebhook(c.ID); err != nil {
		return fmt.Errorf("webhookの削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("Webhookを削除しました (ID: %s)", c.ID))

	return nil
}
