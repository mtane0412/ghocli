/**
 * newsletters.go
 * ニュースレター管理コマンド
 *
 * Ghostニュースレターの管理機能を提供します。
 * Create/Update操作には確認機構が適用されます。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// NewslettersCmd はニュースレター管理コマンドです
type NewslettersCmd struct {
	List   NewslettersListCmd   `cmd:"" help:"List newsletters"`
	Get    NewslettersInfoCmd   `cmd:"" help:"ニュースレターの情報を表示"`
	Create NewslettersCreateCmd `cmd:"" help:"Create a newsletter"`
	Update NewslettersUpdateCmd `cmd:"" help:"Update a newsletter"`
}

// NewslettersListCmd はニュースレター一覧を取得するコマンドです
type NewslettersListCmd struct {
	Limit  int    `help:"Number of newsletters to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)"`
}

// Run はnewslettersコマンドのlistサブコマンドを実行します
func (c *NewslettersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ニュースレター一覧を取得
	response, err := client.ListNewsletters(ghostapi.NewsletterListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
	})
	if err != nil {
		return fmt.Errorf("ニュースレター一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Newsletters)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Name", "Slug", "Status", "Visibility", "Created"}
	rows := make([][]string, len(response.Newsletters))
	for i, newsletter := range response.Newsletters {
		rows[i] = []string{
			newsletter.ID,
			newsletter.Name,
			newsletter.Slug,
			newsletter.Status,
			newsletter.Visibility,
			newsletter.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// NewslettersInfoCmd はニュースレター情報を表示するコマンドです
type NewslettersInfoCmd struct {
	IDOrSlug string `arg:"" help:"Newsletter ID or slug (use 'slug:newsletter-name' format for slug)"`
}

// Run はnewslettersコマンドのinfoサブコマンドを実行します
func (c *NewslettersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ニュースレターを取得
	newsletter, err := client.GetNewsletter(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("ニュースレターの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(newsletter)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", newsletter.ID},
		{"Name", newsletter.Name},
		{"Slug", newsletter.Slug},
		{"Description", newsletter.Description},
		{"Status", newsletter.Status},
		{"Visibility", newsletter.Visibility},
		{"Subscribe on Signup", fmt.Sprintf("%t", newsletter.SubscribeOnSignup)},
		{"Sender Name", newsletter.SenderName},
		{"Sender Email", newsletter.SenderEmail},
		{"Sender Reply To", newsletter.SenderReplyTo},
		{"Sort Order", fmt.Sprintf("%d", newsletter.SortOrder)},
		{"Created", newsletter.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", newsletter.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	return formatter.PrintTable(headers, rows)
}

// NewslettersCreateCmd はニュースレターを作成するコマンドです
type NewslettersCreateCmd struct {
	Name              string `help:"Newsletter name" short:"n" required:""`
	Description       string `help:"Newsletter description" short:"d"`
	Visibility        string `help:"Visibility (members, paid)" default:"members"`
	SubscribeOnSignup bool   `help:"Subscribe members on signup" default:"true"`
	SenderName        string `help:"Sender name"`
	SenderEmail       string `help:"Sender email"`
}

// Run はnewslettersコマンドのcreateサブコマンドを実行します
func (c *NewslettersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("create newsletter '%s'", c.Name)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// 新規ニュースレターを作成
	newNewsletter := &ghostapi.Newsletter{
		Name:              c.Name,
		Description:       c.Description,
		Visibility:        c.Visibility,
		SubscribeOnSignup: c.SubscribeOnSignup,
		SenderName:        c.SenderName,
		SenderEmail:       c.SenderEmail,
	}

	createdNewsletter, err := client.CreateNewsletter(newNewsletter)
	if err != nil {
		return fmt.Errorf("ニュースレターの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ニュースレターを作成しました: %s (ID: %s)", createdNewsletter.Name, createdNewsletter.ID))
	}

	// JSON形式の場合はニュースレター情報も出力
	if root.JSON {
		return formatter.Print(createdNewsletter)
	}

	return nil
}

// NewslettersUpdateCmd はニュースレターを更新するコマンドです
type NewslettersUpdateCmd struct {
	ID                string `arg:"" help:"Newsletter ID"`
	Name              string `help:"Newsletter name" short:"n"`
	Description       string `help:"Newsletter description" short:"d"`
	Visibility        string `help:"Visibility (members, paid)"`
	SubscribeOnSignup *bool  `help:"Subscribe members on signup"`
	SenderName        string `help:"Sender name"`
	SenderEmail       string `help:"Sender email"`
}

// Run はnewslettersコマンドのupdateサブコマンドを実行します
func (c *NewslettersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のニュースレターを取得
	existingNewsletter, err := client.GetNewsletter(c.ID)
	if err != nil {
		return fmt.Errorf("ニュースレターの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("update newsletter '%s' (ID: %s)", existingNewsletter.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// 更新内容を反映
	updateNewsletter := &ghostapi.Newsletter{
		Name:              existingNewsletter.Name,
		Slug:              existingNewsletter.Slug,
		Description:       existingNewsletter.Description,
		Visibility:        existingNewsletter.Visibility,
		SubscribeOnSignup: existingNewsletter.SubscribeOnSignup,
		SenderName:        existingNewsletter.SenderName,
		SenderEmail:       existingNewsletter.SenderEmail,
		SenderReplyTo:     existingNewsletter.SenderReplyTo,
	}

	if c.Name != "" {
		updateNewsletter.Name = c.Name
	}
	if c.Description != "" {
		updateNewsletter.Description = c.Description
	}
	if c.Visibility != "" {
		updateNewsletter.Visibility = c.Visibility
	}
	if c.SubscribeOnSignup != nil {
		updateNewsletter.SubscribeOnSignup = *c.SubscribeOnSignup
	}
	if c.SenderName != "" {
		updateNewsletter.SenderName = c.SenderName
	}
	if c.SenderEmail != "" {
		updateNewsletter.SenderEmail = c.SenderEmail
	}

	// ニュースレターを更新
	updatedNewsletter, err := client.UpdateNewsletter(c.ID, updateNewsletter)
	if err != nil {
		return fmt.Errorf("ニュースレターの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ニュースレターを更新しました: %s (ID: %s)", updatedNewsletter.Name, updatedNewsletter.ID))
	}

	// JSON形式の場合はニュースレター情報も出力
	if root.JSON {
		return formatter.Print(updatedNewsletter)
	}

	return nil
}
