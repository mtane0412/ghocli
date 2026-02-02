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
	"github.com/mtane0412/ghocli/internal/fields"
	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// PagesCmd はページ管理コマンドです
type PagesCmd struct {
	List   PagesListCmd   `cmd:"" help:"List pages"`
	Get    PagesInfoCmd   `cmd:"" help:"Show page information"`
	Cat    PagesCatCmd    `cmd:"" help:"Show content body"`
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
	Copy PagesCopyCmd `cmd:"" help:"Copy a page"`
}

// PagesListCmd is the command to retrieve ページ list
type PagesListCmd struct {
	Status string `help:"Filter by status (draft, published, scheduled, all)" short:"S" default:"all"`
	Limit  int    `help:"Number of pages to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
}

// Run executes the list subcommand of the pages command
func (c *PagesListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PageFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page list
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: c.Status,
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to list pages: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Page struct to map[string]interface{}
		var pagesData []map[string]interface{}
		for _, page := range response.Pages {
			pageMap, err := outfmt.StructToMap(page)
			if err != nil {
				return fmt.Errorf("failed to convert page data: %w", err)
			}
			pagesData = append(pagesData, pageMap)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, pagesData, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// Output in table format
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

// PagesInfoCmd is the command to show ページ information
type PagesInfoCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
}

// Run executes the info subcommand of the pages command
func (c *PagesInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PageFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page
	page, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Page struct to map[string]interface{}
		pageMap, err := outfmt.StructToMap(page)
		if err != nil {
			return fmt.Errorf("failed to convert page data: %w", err)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, []map[string]interface{}{pageMap}, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(page)
	}

	// Output in key/value format (no headers)
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

// PagesCreateCmd is the command to create ページ
type PagesCreateCmd struct {
	Title   string `help:"Page title" short:"t" required:""`
	HTML    string `help:"Page content (HTML)" short:"c"`
	Lexical string `help:"Page content (Lexical JSON)" short:"x"`
	Status  string `help:"Page status (draft, published)" default:"draft"`
}

// Run executes the create subcommand of the pages command
func (c *PagesCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create new page
	newPage := &ghostapi.Page{
		Title:   c.Title,
		HTML:    c.HTML,
		Lexical: c.Lexical,
		Status:  c.Status,
	}

	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created page: %s (ID: %s)", createdPage.Title, createdPage.ID))
	}

	// Also output page information if JSON format
	if root.JSON {
		return formatter.Print(createdPage)
	}

	return nil
}

// PagesUpdateCmd is the command to update ページ
type PagesUpdateCmd struct {
	ID      string `arg:"" help:"Page ID"`
	Title   string `help:"Page title" short:"t"`
	HTML    string `help:"Page content (HTML)" short:"c"`
	Lexical string `help:"Page content (Lexical JSON)" short:"x"`
	Status  string `help:"Page status (draft, published)"`
}

// Run executes the update subcommand of the pages command
func (c *PagesUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing page
	existingPage, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Apply updates
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    existingPage.Status,
		UpdatedAt: existingPage.UpdatedAt, // Use original updated_at from server (for optimistic locking)
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

	// Update page
	updatedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated page: %s (ID: %s)", updatedPage.Title, updatedPage.ID))
	}

	// Also output page information if JSON format
	if root.JSON {
		return formatter.Print(updatedPage)
	}

	return nil
}

// PagesDeleteCmd is the command to delete ページ
type PagesDeleteCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run executes the delete subcommand of the pages command
func (c *PagesDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page information to build confirmation message
	page, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete page '%s' (ID: %s)", page.Title, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Delete page
	if err := client.DeletePage(c.ID); err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted page (ID: %s)", c.ID))

	return nil
}

// ========================================
// Phase 1: URL取得
// ========================================

// PagesURLCmd is the command to get ページ web URL
type PagesURLCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
	Open     bool   `help:"Open URL in browser" short:"o"`
}

// Run executes the url subcommand of the pages command
func (c *PagesURLCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page
	page, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Get URL
	url := page.URL
	if url == "" {
		return fmt.Errorf("could not get page URL")
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output URL
	formatter.PrintMessage(url)

	// Open in browser if --open flag is specified
	if c.Open {
		// Open browser with OS-appropriate command
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
			return fmt.Errorf("failed to open URL in browser: %w", err)
		}
	}

	return nil
}

// ========================================
// Phase 2: 状態変更
// ========================================

// PagesPublishCmd is the command to publish ページ
type PagesPublishCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run executes the publish subcommand of the pages command
func (c *PagesPublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing page
	existingPage, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Error if already published
	if existingPage.Status == "published" {
		return fmt.Errorf("page is already published")
	}

	// Change status to published
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    "published",
		UpdatedAt: existingPage.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update page
	publishedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("failed to publish page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("published page: %s (ID: %s)", publishedPage.Title, publishedPage.ID))
	}

	// Also output page information if JSON format
	if root.JSON {
		return formatter.Print(publishedPage)
	}

	return nil
}

// PagesUnpublishCmd is the command to unpublish ページ
type PagesUnpublishCmd struct {
	ID string `arg:"" help:"Page ID"`
}

// Run executes the unpublish subcommand of the pages command
func (c *PagesUnpublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing page
	existingPage, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Error if already a draft
	if existingPage.Status == "draft" {
		return fmt.Errorf("page is already a draft")
	}

	// Change status to draft
	updatePage := &ghostapi.Page{
		Title:     existingPage.Title,
		Slug:      existingPage.Slug,
		HTML:      existingPage.HTML,
		Lexical:   existingPage.Lexical,
		Status:    "draft",
		UpdatedAt: existingPage.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update page
	unpublishedPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("failed to unpublish page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("unpublished page: %s (ID: %s)", unpublishedPage.Title, unpublishedPage.ID))
	}

	// Also output page information if JSON format
	if root.JSON {
		return formatter.Print(unpublishedPage)
	}

	return nil
}

// ========================================
// Phase 2: catコマンド
// ========================================

// PagesCatCmd is the command to show ページ content body
type PagesCatCmd struct {
	IDOrSlug string `arg:"" help:"Page ID or slug"`
	Format   string `help:"Output format (text, html, lexical)" default:"text"`
}

// Run executes the cat subcommand of the pages command
func (c *PagesCatCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page
	page, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output according to format
	var content string
	switch c.Format {
	case "html":
		content = page.HTML
	case "text":
		// Convert HTML to text
		content = html2text.HTML2Text(page.HTML)
	case "lexical":
		content = page.Lexical
	default:
		return fmt.Errorf("unsupported format: %s (html, text, lexical のいずれかを指定してください)", c.Format)
	}

	// Output content
	formatter.PrintMessage(content)

	return nil
}

// ========================================
// Phase 8.3: copyコマンド
// ========================================

// PagesCopyCmd is the command to copy ページ
type PagesCopyCmd struct {
	IDOrSlug string `arg:"" help:"Source page ID or slug"`
	Title    string `help:"New title (defaults to 'Original Title (Copy)')" short:"t"`
}

// Run executes the copy subcommand of the pages command
func (c *PagesCopyCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get original page
	original, err := client.GetPage(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Determine new title
	newTitle := c.Title
	if newTitle == "" {
		newTitle = original.Title + " (Copy)"
	}

	// Create new page (exclude ID/UUID/Slug/URL/dates, Status fixed to draft)
	newPage := &ghostapi.Page{
		Title:   newTitle,
		HTML:    original.HTML,
		Lexical: original.Lexical,
		Status:  "draft",
	}

	// Create page
	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		return fmt.Errorf("failed to copy page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("copied page: %s (ID: %s)", createdPage.Title, createdPage.ID))
	}

	// Also output page information if JSON format
	if root.JSON {
		return formatter.Print(createdPage)
	}

	return nil
}

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// PagesDraftsCmd is the command to retrieve draft ページ list
type PagesDraftsCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the drafts subcommand of the pages command
func (c *PagesDraftsCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get draft page list
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "draft",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("下書きfailed to list pages: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// Output in table format
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

// PagesPublishedCmd is the command to retrieve published ページ list
type PagesPublishedCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the published subcommand of the pages command
func (c *PagesPublishedCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get published page list
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "published",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("公開済みfailed to list pages: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// Output in table format
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

// PagesScheduledCmd is the command to retrieve scheduled ページ list
type PagesScheduledCmd struct {
	Limit int `help:"Number of pages to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the scheduled subcommand of the pages command
func (c *PagesScheduledCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get scheduled page list
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "scheduled",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("予約failed to list pages: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Pages)
	}

	// Output in table format
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

// PagesScheduleCmd is the command to schedule ページ for publishing
type PagesScheduleCmd struct {
	ID string `arg:"" help:"Page ID"`
	At string `help:"Schedule time (YYYY-MM-DD HH:MM)" required:""`
}

// Run executes the schedule subcommand of the pages command
func (c *PagesScheduleCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing page
	existingPage, err := client.GetPage(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	// Parse datetime
	publishedAt, err := parseDateTime(c.At)
	if err != nil {
		return fmt.Errorf("failed to parse datetime: %w", err)
	}

	// Change status to scheduled and set publish date
	updatePage := &ghostapi.Page{
		Title:       existingPage.Title,
		Slug:        existingPage.Slug,
		HTML:        existingPage.HTML,
		Lexical:     existingPage.Lexical,
		Status:      "scheduled",
		PublishedAt: &publishedAt,
		UpdatedAt:   existingPage.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update page
	scheduledPage, err := client.UpdatePage(c.ID, updatePage)
	if err != nil {
		return fmt.Errorf("failed to schedule page: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("scheduled page: %s (ID: %s, scheduled for: %s)",
			scheduledPage.Title, scheduledPage.ID, publishedAt.Format("2006-01-02 15:04")))
	}

	// Also output page information if JSON format
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

// PagesBatchPublishCmd is the command to batch publish ページ
type PagesBatchPublishCmd struct {
	IDs []string `arg:"" help:"Page IDs to publish"`
}

// Run executes the pages batch publish subcommand
func (c *PagesBatchPublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Publish each page
	successCount := 0
	for _, id := range c.IDs {
		// Get existing page
		existingPage, err := client.GetPage(id)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to get page (ID: %s): %v", id, err))
			continue
		}

		// すでに公開済みの場合はスキップ
		if existingPage.Status == "published" {
			formatter.PrintMessage(fmt.Sprintf("skipped (already published): %s (ID: %s)", existingPage.Title, id))
			continue
		}

		// Change status to published
		updatePage := &ghostapi.Page{
			Title:     existingPage.Title,
			Slug:      existingPage.Slug,
			HTML:      existingPage.HTML,
			Lexical:   existingPage.Lexical,
			Status:    "published",
			UpdatedAt: existingPage.UpdatedAt,
		}

		// Update page
		_, err = client.UpdatePage(id, updatePage)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to publish page (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("published: %s (ID: %s)", existingPage.Title, id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\ncompleted: %d件のpublished page", successCount))

	return nil
}

// PagesBatchDeleteCmd is the command to batch delete ページ
type PagesBatchDeleteCmd struct {
	IDs []string `arg:"" help:"Page IDs to delete"`
}

// Run executes the pages batch delete subcommand
func (c *PagesBatchDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete %d pages", len(c.IDs))
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Delete each page
	successCount := 0
	for _, id := range c.IDs {
		// Delete page
		if err := client.DeletePage(id); err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to delete page (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("deleted (ID: %s)", id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\ncompleted: %d件のdeleted page", successCount))

	return nil
}

// ========================================
// Phase 4: ページ検索
// ========================================

// PagesSearchCmd is the command to search ページ
type PagesSearchCmd struct {
	Query string `arg:"" help:"Search query"`
	Limit int    `help:"Number of pages to retrieve" short:"l" aliases:"max,n" default:"15"`
}

// Run executes the search subcommand of the pages command
func (c *PagesSearchCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get page list（検索クエリはfilterとして渡す）
	response, err := client.ListPages(ghostapi.ListOptions{
		Status: "all",
		Limit:  c.Limit,
		Page:   1,
	})
	if err != nil {
		return fmt.Errorf("failed to search pages: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter pages matching query (simple implementation)
	var filteredPages []ghostapi.Page
	for _, page := range response.Pages {
		if containsIgnoreCase(page.Title, c.Query) || containsIgnoreCase(page.HTML, c.Query) {
			filteredPages = append(filteredPages, page)
		}
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(filteredPages)
	}

	// Output in table format
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
