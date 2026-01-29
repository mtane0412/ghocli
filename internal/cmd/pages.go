/**
 * pages.go
 * ページ管理コマンド
 *
 * Ghostページの作成、更新、削除機能を提供します。
 */

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// PagesCmd はページ管理コマンドです
type PagesCmd struct {
	List   PagesListCmd   `cmd:"" help:"List pages"`
	Get    PagesGetCmd    `cmd:"" help:"Get a page"`
	Create PagesCreateCmd `cmd:"" help:"Create a page"`
	Update PagesUpdateCmd `cmd:"" help:"Update a page"`
	Delete PagesDeleteCmd `cmd:"" help:"Delete a page"`
}

// PagesListCmd はページ一覧を取得するコマンドです
type PagesListCmd struct {
	Status string `help:"Filter by status (draft, published, scheduled, all)" short:"S" default:"all"`
	Limit  int    `help:"Number of pages to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
}

// Run はpagesコマンドのlistサブコマンドを実行します
func (c *PagesListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ページ一覧を取得
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: c.Status,
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("ページ一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Title", "Status", "Created", "Published"}
	rows := make([][]string, len(response.Pages))
	for i, page := range response.Pages {
		publishedAt := ""
		if page.PublishedAt != nil {
			publishedAt = page.PublishedAt.Format("2006-01-02")
		}
		rows[i] = []string{
			page.ID,
			page.Title,
			page.Status,
			page.CreatedAt.Format("2006-01-02"),
			publishedAt,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// PagesGetCmd はページを取得するコマンドです
type PagesGetCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
}

// Run はpagesコマンドのgetサブコマンドを実行します
func (c *PagesGetCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ページを取得
	page, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("ページの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(page)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", page.ID},
		{"Title", page.Title},
		{"Slug", page.Slug},
		{"Status", page.Status},
		{"Created", page.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", page.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if page.PublishedAt != nil {
		rows = append(rows, []string{"Published", page.PublishedAt.Format("2006-01-02 15:04:05")})
	}

	return formatter.PrintTable(headers, rows)
}

// PagesCreateCmd はページを作成するコマンドです
type PagesCreateCmd struct {
	Title   string `help:"Page title" short:"t" required:""`
	HTML    string `help:"Page content (HTML)" short:"c"`
	Lexical string `help:"Page content (Lexical JSON)" short:"x"`
	Status  string `help:"Page status (draft, published)" default:"draft"`
}

// Run はpagesコマンドのcreateサブコマンドを実行します
func (c *PagesCreateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 新規ページを作成
	newPage := &ghostapi.Page{
		Title:   c.Title,
		HTML:    c.HTML,
		Lexical: c.Lexical,
		Status:  c.Status,
	}

	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		return fmt.Errorf("ページの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページを作成しました: %s (ID: %s)", createdPage.Title, createdPage.ID))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(createdPage)
	}

	return nil
}

// PagesUpdateCmd はページを更新するコマンドです
type PagesUpdateCmd struct {
	ID      string `arg:"" help:"Page ID"`
	Title   string `help:"Page title" short:"t"`
	HTML    string `help:"Page content (HTML)" short:"c"`
	Lexical string `help:"Page content (Lexical JSON)" short:"x"`
	Status  string `help:"Page status (draft, published)"`
}

// Run はpagesコマンドのupdateサブコマンドを実行します
func (c *PagesUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のページを取得
	existingPage, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("ページの取得に失敗: %w", err)
	}

	// 更新内容を反映
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    existingPage.Status,
		UpdatedAt: time.Now(), // 更新時刻を設定
	}

	if c.Title != "" {
		updatePage.Title = c.Title
	}
	if c.HTML != "" {
		updatePage.HTML = c.HTML
	}
	if c.Lexical != "" {
		updatePage.Lexical = c.Lexical
	}
	if c.Status != "" {
		updatePage.Status = c.Status
	}

	// ページを更新
	updatedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("ページの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページを更新しました: %s (ID: %s)", updatedPage.Title, updatedPage.ID))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(updatedPage)
	}

	return nil
}

// PagesDeleteCmd はページを削除するコマンドです
type PagesDeleteCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run はpagesコマンドのdeleteサブコマンドを実行します
func (c *PagesDeleteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 確認なしで削除する場合を除き、確認を求める
	if !root.Force {
		// ページ情報を取得して確認
		page, err := client.GetPage(c.ID)
		if err != nil {
			return fmt.Errorf("ページの取得に失敗: %w", err)
		}

		fmt.Printf("本当にページ「%s」(ID: %s)を削除しますか? [y/N]: ", page.Title, c.ID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("削除がキャンセルされました")
		}
	}

	// ページを削除
	if err := client.DeletePage(c.ID); err != nil {
		return fmt.Errorf("ページの削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("ページを削除しました (ID: %s)", c.ID))

	return nil
}
