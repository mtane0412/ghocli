/**
 * newsletters.go
 * ニュースレター管理コマンド
 *
 * Ghostニュースレターの閲覧機能を提供します。
 * ビジネス設定の誤変更リスクを回避するため、読み取り操作（List, Get）のみ実装しています。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// NewslettersCmd はニュースレター管理コマンドです
type NewslettersCmd struct {
	List NewslettersListCmd `cmd:"" help:"List newsletters"`
	Get  NewslettersGetCmd  `cmd:"" help:"Get a newsletter"`
}

// NewslettersListCmd はニュースレター一覧を取得するコマンドです
type NewslettersListCmd struct {
	Limit  int    `help:"Number of newsletters to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)"`
}

// Run はnewslettersコマンドのlistサブコマンドを実行します
func (c *NewslettersListCmd) Run(root *RootFlags) error {
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

// NewslettersGetCmd はニュースレターを取得するコマンドです
type NewslettersGetCmd struct {
	IDOrSlug string `arg:"" help:"Newsletter ID or slug (use 'slug:newsletter-name' format for slug)"`
}

// Run はnewslettersコマンドのgetサブコマンドを実行します
func (c *NewslettersGetCmd) Run(root *RootFlags) error {
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
