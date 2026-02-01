/**
 * pages.go
 * ページ管理コマンド
 *
 * Ghostページの作成、更新、削除機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/k3a/html2text"
	"github.com/mtane0412/gho/internal/fields"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// PagesCmd はページ管理コマンドです
type PagesCmd struct {
	List   PagesListCmd   `cmd:"" help:"List pages"`
	Get    PagesInfoCmd   `cmd:"" help:"ページの情報を表示"`
	Cat    PagesCatCmd    `cmd:"" help:"本文コンテンツを表示"`
	Create PagesCreateCmd `cmd:"" help:"Create a page"`
	Update PagesUpdateCmd `cmd:"" help:"Update a page"`
	Delete PagesDeleteCmd `cmd:"" help:"Delete a page"`

	// Phase 1: ステータス別一覧ショートカット
	Drafts    PagesDraftsCmd    `cmd:"" help:"List draft pages"`
	Published PagesPublishedCmd `cmd:"" help:"List published pages"`
	Scheduled PagesScheduledCmd `cmd:"" help:"List scheduled pages"`

	// Phase 1: URL取得
	URL PagesURLCmd `cmd:"" help:"Get page URL"`

	// Phase 2: 状態変更
	Publish   PagesPublishCmd   `cmd:"" help:"Publish a page"`
	Unpublish PagesUnpublishCmd `cmd:"" help:"Unpublish a page"`

	// Phase 3: 予約公開
	Schedule PagesScheduleCmd `cmd:"" help:"Schedule a page"`

	// Phase 4: バッチ操作
	Batch  PagesBatchCmd  `cmd:"" help:"Batch operations"`
	Search PagesSearchCmd `cmd:"" help:"Search pages"`

	// Phase 8.3: コピー
	Copy PagesCopyCmd `cmd:"" help:"ページをコピー"`
}

// PagesListCmd はページ一覧を取得するコマンドです
type PagesListCmd struct {
	Status string `help:"Filter by status (draft, published, scheduled, all)" short:"S" default:"all"`
	Limit  int    `help:"Number of pages to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
}

// Run はpagesコマンドのlistサブコマンドを実行します
func (c *PagesListCmd) Run(ctx context.Context, root *RootFlags) error {
	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PageFields)
		if err != nil {
			return fmt.Errorf("フィールド指定のパースに失敗: %w", err)
		}
		selectedFields = parsedFields
	}

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

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Page構造体をmap[string]interface{}に変換
		var pagesData []map[string]interface{}
		for _, page := range response.Pages {
			pageMap, err := outfmt.StructToMap(page)
			if err != nil {
				return fmt.Errorf("ページデータの変換に失敗: %w", err)
			}
			pagesData = append(pagesData, pageMap)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, pagesData, selectedFields)
	}

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

// PagesInfoCmd はページ情報を表示するコマンドです
type PagesInfoCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
}

// Run はpagesコマンドのinfoサブコマンドを実行します
func (c *PagesInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PageFields)
		if err != nil {
			return fmt.Errorf("フィールド指定のパースに失敗: %w", err)
		}
		selectedFields = parsedFields
	}

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

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Page構造体をmap[string]interface{}に変換
		pageMap, err := outfmt.StructToMap(page)
		if err != nil {
			return fmt.Errorf("ページデータの変換に失敗: %w", err)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, []map[string]interface{}{pageMap}, selectedFields)
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(page)
	}

	// キー/値形式で出力（ヘッダーなし）
	rows := [][]string{
		{"id", page.ID},
		{"title", page.Title},
		{"slug", page.Slug},
		{"status", page.Status},
	}

	// visibilityを追加
	if page.Visibility != "" {
		rows = append(rows, []string{"visibility", page.Visibility})
	}

	// urlを追加
	if page.URL != "" {
		rows = append(rows, []string{"url", page.URL})
	}

	// authorsを追加（存在する場合のみ）
	if len(page.Authors) > 0 {
		rows = append(rows, []string{"authors", outfmt.FormatAuthors(page.Authors)})
	}

	// tagsを追加（存在する場合のみ）
	if len(page.Tags) > 0 {
		rows = append(rows, []string{"tags", outfmt.FormatTags(page.Tags)})
	}

	// featuredを追加（trueの場合のみ）
	if page.Featured {
		rows = append(rows, []string{"featured", "true"})
	}

	// excerptを追加（存在する場合のみ、140文字で切り詰め）
	if page.Excerpt != "" {
		rows = append(rows, []string{"excerpt", outfmt.TruncateExcerpt(page.Excerpt, 140)})
	}

	// 日時フィールド
	rows = append(rows, []string{"created", page.CreatedAt.Format("2006-01-02 15:04:05")})
	rows = append(rows, []string{"updated", page.UpdatedAt.Format("2006-01-02 15:04:05")})

	if page.PublishedAt != nil {
		rows = append(rows, []string{"published", page.PublishedAt.Format("2006-01-02 15:04:05")})
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// PagesCreateCmd はページを作成するコマンドです
type PagesCreateCmd struct {
	Title   string `help:"Page title" short:"t" required:""`
	HTML    string `help:"Page content (HTML)" short:"c"`
	Lexical string `help:"Page content (Lexical JSON)" short:"x"`
	Status  string `help:"Page status (draft, published)" default:"draft"`
}

// Run はpagesコマンドのcreateサブコマンドを実行します
func (c *PagesCreateCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *PagesUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *PagesDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
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
	if err := ConfirmDestructive(ctx, root, action); err != nil {
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
func (c *PagesURLCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *PagesPublishCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *PagesUnpublishCmd) Run(ctx context.Context, root *RootFlags) error {
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

// ========================================
// Phase 2: catコマンド
// ========================================

// PagesCatCmd はページの本文コンテンツを表示するコマンドです
type PagesCatCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
	Format   string `help:"Output format (text, html, lexical)" default:"text"`
}

// Run はpagesコマンドのcatサブコマンドを実行します
func (c *PagesCatCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// フォーマットに応じて出力
	var content string
	switch c.Format {
	case "html":
		content = page.HTML
	case "text":
		// HTMLからテキストへ変換
		content = html2text.HTML2Text(page.HTML)
	case "lexical":
		content = page.Lexical
	default:
		return fmt.Errorf("未対応のフォーマット: %s (html, text, lexical のいずれかを指定してください)", c.Format)
	}

	// コンテンツを出力
	formatter.PrintMessage(content)

	return nil
}

// ========================================
// Phase 8.3: copyコマンド
// ========================================

// PagesCopyCmd はページをコピーするコマンドです
type PagesCopyCmd struct {
	IDOrSlug string `arg:"" help:"コピー元のページID またはスラッグ"`
	Title    string `help:"新しいタイトル（省略時は '元タイトル (Copy)'）" short:"t"`
}

// Run はpagesコマンドのcopyサブコマンドを実行します
func (c *PagesCopyCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 元のページを取得
	original, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("ページの取得に失敗: %w", err)
	}

	// 新しいタイトルを決定
	newTitle := c.Title
	if newTitle == "" {
		newTitle = original.Title + " (Copy)"
	}

	// 新しいページを作成（ID/UUID/Slug/URL/日時は除外、Statusはdraft固定）
	newPage := &ghostapi.Page{
		Title:   newTitle,
		HTML:    original.HTML,
		Lexical: original.Lexical,
		Status:  "draft",
	}

	// ページを作成
	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		return fmt.Errorf("ページのコピーに失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページをコピーしました: %s (ID: %s)", createdPage.Title, createdPage.ID))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(createdPage)
	}

	return nil
}

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// PagesDraftsCmd は下書きページ一覧を取得するコマンドです
type PagesDraftsCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpagesコマンドのdraftsサブコマンドを実行します
func (c *PagesDraftsCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 下書きページ一覧を取得
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "draft",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("下書きページ一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Title", "Status", "Created", "Updated"}
	rows := make([][]string, len(response.Pages))
	for i, page := range response.Pages {
		rows[i] = []string{
			page.ID,
			page.Title,
			page.Status,
			page.CreatedAt.Format("2006-01-02"),
			page.UpdatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// PagesPublishedCmd は公開済みページ一覧を取得するコマンドです
type PagesPublishedCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpagesコマンドのpublishedサブコマンドを実行します
func (c *PagesPublishedCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 公開済みページ一覧を取得
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "published",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("公開済みページ一覧の取得に失敗: %w", err)
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

// PagesScheduledCmd は予約ページ一覧を取得するコマンドです
type PagesScheduledCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpagesコマンドのscheduledサブコマンドを実行します
func (c *PagesScheduledCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 予約ページ一覧を取得
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "scheduled",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("予約ページ一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Title", "Status", "Created", "Scheduled"}
	rows := make([][]string, len(response.Pages))
	for i, page := range response.Pages {
		publishedAt := ""
		if page.PublishedAt != nil {
			publishedAt = page.PublishedAt.Format("2006-01-02 15:04")
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

// ========================================
// Phase 3: 予約公開
// ========================================

// PagesScheduleCmd はページを予約公開に設定するコマンドです
type PagesScheduleCmd struct {
	ID string `arg:"" help:"Page ID"`
	At string `help:"Schedule time (YYYY-MM-DD HH:MM)" required:""`
}

// Run はpagesコマンドのscheduleサブコマンドを実行します
func (c *PagesScheduleCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// 日時をパース
	publishedAt, err := parseDateTime(c.At)
	if err != nil {
		return fmt.Errorf("日時のパースに失敗: %w", err)
	}

	// ステータスをscheduledに変更し、公開日時を設定
	updatePage := &ghostapi.Page{
		Title:       existingPage.Title,
		Slug:        existingPage.Slug,
		HTML:        existingPage.HTML,
		Lexical:     existingPage.Lexical,
		Status:      "scheduled",
		PublishedAt: &publishedAt,
		UpdatedAt:   existingPage.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// ページを更新
	scheduledPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("ページの予約公開設定に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ページを予約公開に設定しました: %s (ID: %s, 公開予定: %s)",
			scheduledPage.Title, scheduledPage.ID, publishedAt.Format("2006-01-02 15:04")))
	}

	// JSON形式の場合はページ情報も出力
	if root.JSON {
		return formatter.Print(scheduledPage)
	}

	return nil
}

// ========================================
// Phase 4: バッチ操作
// ========================================

// PagesBatchCmd はバッチ操作コマンドです
type PagesBatchCmd struct {
	Publish PagesBatchPublishCmd `cmd:"" help:"Batch publish pages"`
	Delete  PagesBatchDeleteCmd  `cmd:"" help:"Batch delete pages"`
}

// PagesBatchPublishCmd は複数ページを一括公開するコマンドです
type PagesBatchPublishCmd struct {
	IDs []string `arg:"" help:"Page IDs to publish"`
}

// Run はpages batch publishサブコマンドを実行します
func (c *PagesBatchPublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 各ページを公開
	successCount := 0
	for _, id := range c.IDs {
		// 既存のページを取得
		existingPage, err := client.GetPage(id)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("ページの取得に失敗 (ID: %s): %v", id, err))
			continue
		}

		// すでに公開済みの場合はスキップ
		if existingPage.Status == "published" {
			formatter.PrintMessage(fmt.Sprintf("スキップ (すでに公開済み): %s (ID: %s)", existingPage.Title, id))
			continue
		}

		// ステータスをpublishedに変更
		updatePage := &ghostapi.Page{
			Title:     existingPage.Title,
			Slug:      existingPage.Slug,
			HTML:      existingPage.HTML,
			Lexical:   existingPage.Lexical,
			Status:    "published",
			UpdatedAt: existingPage.UpdatedAt,
		}

		// ページを更新
		_, err = client.UpdatePage(id, updatePage)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("ページの公開に失敗 (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("公開しました: %s (ID: %s)", existingPage.Title, id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\n完了: %d件のページを公開しました", successCount))

	return nil
}

// PagesBatchDeleteCmd は複数ページを一括削除するコマンドです
type PagesBatchDeleteCmd struct {
	IDs []string `arg:"" help:"Page IDs to delete"`
}

// Run はpages batch deleteサブコマンドを実行します
func (c *PagesBatchDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete %d pages", len(c.IDs))
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 各ページを削除
	successCount := 0
	for _, id := range c.IDs {
		// ページを削除
		if err := client.DeletePage(id); err != nil {
			formatter.PrintMessage(fmt.Sprintf("ページの削除に失敗 (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("削除しました (ID: %s)", id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\n完了: %d件のページを削除しました", successCount))

	return nil
}

// ========================================
// Phase 4: ページ検索
// ========================================

// PagesSearchCmd はページを検索するコマンドです
type PagesSearchCmd struct {
	Query string `arg:"" help:"Search query"`
	Limit int    `help:"Number of pages to retrieve" short:"l" default:"15"`
}

// Run はpagesコマンドのsearchサブコマンドを実行します
func (c *PagesSearchCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ページ一覧を取得（検索クエリはfilterとして渡す）
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "all",
		Limit:  c.Limit,
		Page:   1,
	})
	if err != nil {
		return fmt.Errorf("ページ検索に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// クエリに一致するページをフィルタリング（簡易的な実装）
	var filteredPages []ghostapi.Page
	for _, page := range response.Pages {
		if containsIgnoreCase(page.Title, c.Query) || containsIgnoreCase(page.HTML, c.Query) {
			filteredPages = append(filteredPages, page)
		}
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(filteredPages)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Title", "Status", "Created"}
	rows := make([][]string, len(filteredPages))
	for i, page := range filteredPages {
		rows[i] = []string{
			page.ID,
			page.Title,
			page.Status,
			page.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}
