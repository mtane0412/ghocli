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

	// Phase 1: URL取得
	URL PagesURLCmd `cmd:"" help:"Get page URL"`

	// Phase 2: 状態変更
	Publish   PagesPublishCmd   `cmd:"" help:"Publish a page"`
	Unpublish PagesUnpublishCmd `cmd:"" help:"Unpublish a page"`
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
		UpdatedAt: existingPage.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
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

	// ページ情報を取得して確認メッセージを構築
	page, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("ページの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete page '%s' (ID: %s)", page.Title, c.ID)
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
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

// ========================================
// Phase 1: URL取得
// ========================================

// PagesURLCmd はページのWeb URLを取得するコマンドです
type PagesURLCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
	Open     bool   `help:"Open URL in browser" short:"o"`
}

// Run はpagesコマンドのurlサブコマンドを実行します
func (c *PagesURLCmd) Run(root *RootFlags) error {
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

	// URLを取得
	url := page.URL
	if url == "" {
		return fmt.Errorf("ページのURLが取得できませんでした")
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// URLを出力
	formatter.PrintMessage(url)

	// --openフラグが指定されている場合はブラウザで開く
	if c.Open {
		// OSに応じたコマンドでブラウザを開く
		var cmd string
		switch {
		case fileExists("/usr/bin/open"): // macOS
			cmd = "open"
		case fileExists("/usr/bin/xdg-open"): // Linux
			cmd = "xdg-open"
		default: // Windows
			cmd = "start"
		}

		if err := runCommand(cmd, url); err != nil {
			return fmt.Errorf("ブラウザでURLを開くことに失敗: %w", err)
		}
	}

	return nil
}

// ========================================
// Phase 2: 状態変更
// ========================================

// PagesPublishCmd はページを公開するコマンドです
type PagesPublishCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run はpagesコマンドのpublishサブコマンドを実行します
func (c *PagesPublishCmd) Run(root *RootFlags) error {
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

	// すでに公開済みの場合はエラー
	if existingPage.Status == "published" {
		return fmt.Errorf("このページはすでに公開されています")
	}

	// ステータスをpublishedに変更
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    "published",
		UpdatedAt: existingPage.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// ページを更新
	publishedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("ページの公開に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページを公開しました: %s (ID: %s)", publishedPage.Title, publishedPage.ID))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(publishedPage)
	}

	return nil
}

// PagesUnpublishCmd はページを下書きに戻すコマンドです
type PagesUnpublishCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run はpagesコマンドのunpublishサブコマンドを実行します
func (c *PagesUnpublishCmd) Run(root *RootFlags) error {
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

	// すでに下書きの場合はエラー
	if existingPage.Status == "draft" {
		return fmt.Errorf("このページはすでに下書きです")
	}

	// ステータスをdraftに変更
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    "draft",
		UpdatedAt: existingPage.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// ページを更新
	unpublishedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("ページの非公開化に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページを下書きに戻しました: %s (ID: %s)", unpublishedPage.Title, unpublishedPage.ID))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(unpublishedPage)
	}

	return nil
}
