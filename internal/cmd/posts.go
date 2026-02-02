/**
 * posts.go
 * Post management commands
 *
 * Provides functionality for creating, updating, deleting, and publishing Ghost posts.
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/k3a/html2text"
	"github.com/mtane0412/ghocli/internal/fields"
	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/input"
	"github.com/mtane0412/ghocli/internal/markdown"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// PostsCmd is the post management command
type PostsCmd struct {
	List    PostsListCmd    `cmd:"" help:"List posts"`
	Get     PostsInfoCmd    `cmd:"" help:"Show post information"`
	Cat     PostsCatCmd     `cmd:"" help:"Show content body"`
	Create  PostsCreateCmd  `cmd:"" help:"Create a post"`
	Update  PostsUpdateCmd  `cmd:"" help:"Update a post"`
	Delete  PostsDeleteCmd  `cmd:"" help:"Delete a post"`
	Publish PostsPublishCmd `cmd:"" help:"Publish a draft"`

	// Phase 1: Status-based list shortcuts
	Drafts    PostsDraftsCmd    `cmd:"" help:"List draft posts"`
	Published PostsPublishedCmd `cmd:"" help:"List published posts"`
	Scheduled PostsScheduledCmd `cmd:"" help:"List scheduled posts"`

	// Phase 1: URL retrieval
	URL PostsURLCmd `cmd:"" help:"Get post URL"`

	// Phase 2: State changes
	Unpublish PostsUnpublishCmd `cmd:"" help:"Unpublish a post"`

	// Phase 3: Scheduled publishing
	Schedule PostsScheduleCmd `cmd:"" help:"Schedule a post"`

	// Phase 4: Batch operations
	Batch  PostsBatchCmd  `cmd:"" help:"Batch operations"`
	Search PostsSearchCmd `cmd:"" help:"Search posts"`

	// Phase 8.3: Copy
	Copy PostsCopyCmd `cmd:"" help:"Copy a post"`
}

// PostsListCmd is the command to retrieve 投稿 list
type PostsListCmd struct {
	Status string `help:"Filter by status (draft, published, scheduled, all)" short:"S" default:"all"`
	Limit  int    `help:"Number of posts to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
}

// Run executes the list subcommand of the posts command
func (c *PostsListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PostFields)
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

	// Get post list
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: c.Status,
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to list posts: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Post struct to map[string]interface{}
		var postsData []map[string]interface{}
		for _, post := range response.Posts {
			postMap, err := outfmt.StructToMap(post)
			if err != nil {
				return fmt.Errorf("failed to convert post data: %w", err)
			}
			postsData = append(postsData, postMap)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, postsData, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// Output in table format
	headers := []string{"ID", "Title", "Status", "Created", "Published"}
	rows := make([][]string, len(response.Posts))
	for i, post := range response.Posts {
		publishedAt := ""
		if post.PublishedAt != nil {
			publishedAt = post.PublishedAt.Format("2006-01-02")
		}
		rows[i] = []string{
			post.ID,
			post.Title,
			post.Status,
			post.CreatedAt.Format("2006-01-02"),
			publishedAt,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// PostsInfoCmd is the command to show 投稿 information
type PostsInfoCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
}

// Run executes the info subcommand of the posts command
func (c *PostsInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.PostFields)
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

	// Get post
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Post struct to map[string]interface{}
		postMap, err := outfmt.StructToMap(post)
		if err != nil {
			return fmt.Errorf("failed to convert post data: %w", err)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, []map[string]interface{}{postMap}, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(post)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"id", post.ID},
		{"title", post.Title},
		{"slug", post.Slug},
		{"status", post.Status},
	}

	// visibilityを追加
	if post.Visibility != "" {
		rows = append(rows, []string{"visibility", post.Visibility})
	}

	// urlを追加
	if post.URL != "" {
		rows = append(rows, []string{"url", post.URL})
	}

	// authorsを追加（存在する場合のみ）
	if len(post.Authors) > 0 {
		rows = append(rows, []string{"authors", outfmt.FormatAuthors(post.Authors)})
	}

	// tagsを追加（存在する場合のみ）
	if len(post.Tags) > 0 {
		rows = append(rows, []string{"tags", outfmt.FormatTags(post.Tags)})
	}

	// featuredを追加（trueの場合のみ）
	if post.Featured {
		rows = append(rows, []string{"featured", "true"})
	}

	// excerptを追加（存在する場合のみ、140文字で切り詰め）
	if post.Excerpt != "" {
		rows = append(rows, []string{"excerpt", outfmt.TruncateExcerpt(post.Excerpt, 140)})
	}

	// 日時フィールド
	rows = append(rows, []string{"created", post.CreatedAt.Format("2006-01-02 15:04:05")})
	rows = append(rows, []string{"updated", post.UpdatedAt.Format("2006-01-02 15:04:05")})

	if post.PublishedAt != nil {
		rows = append(rows, []string{"published", post.PublishedAt.Format("2006-01-02 15:04:05")})
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// PostsCreateCmd is the command to create 投稿
type PostsCreateCmd struct {
	Title    string `help:"Post title" short:"t" required:""`
	HTML     string `help:"Post content (HTML)" short:"c"`
	Markdown string `help:"Post content (Markdown)" short:"m"`
	Lexical  string `help:"Post content (Lexical JSON)" short:"x"`
	File     string `help:"Read content from file (auto-detect format)" type:"existingfile"`
	Status   string `help:"Post status (draft, published)" default:"draft"`
}

// Run executes the create subcommand of the posts command
func (c *PostsCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// コンテンツとフォーマットの決定
	var htmlContent string
	var format input.ContentFormat

	// ファイル指定の場合はフォーマット自動検出
	if c.File != "" {
		fileContent, detectedFormat, err := input.ReadContentWithFormat(c.File, "")
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		format = detectedFormat

		// フォーマットに応じて処理
		switch format {
		case input.FormatMarkdown:
			// Markdown→HTML変換
			htmlContent, err = markdown.ConvertToHTML(fileContent)
			if err != nil {
				return fmt.Errorf("failed to convert markdown to HTML: %w", err)
			}
		case input.FormatHTML:
			// HTMLはそのまま使用
			htmlContent = fileContent
		case input.FormatLexical:
			// Lexical JSONはそのまま使用（c.Lexicalに設定）
			c.Lexical = fileContent
		default:
			// 不明な形式の場合はHTMLとして扱う
			htmlContent = fileContent
		}
	} else {
		// インラインコンテンツの処理
		htmlContent = c.HTML

		// Markdownフラグが指定されている場合はMarkdown→HTML変換
		if c.Markdown != "" {
			htmlContent, err = markdown.ConvertToHTML(c.Markdown)
			if err != nil {
				return fmt.Errorf("failed to convert markdown to HTML: %w", err)
			}
		}
	}

	// Create new post
	newPost := &ghostapi.Post{
		Title:   c.Title,
		HTML:    htmlContent,
		Lexical: c.Lexical,
		Status:  c.Status,
	}

	// HTMLコンテンツが指定されている場合は、自動的にsource=htmlを適用
	var createdPost *ghostapi.Post
	if htmlContent != "" && c.Lexical == "" {
		// HTMLをサーバー側でLexical形式に変換
		opts := ghostapi.CreateOptions{
			Source: "html",
		}
		createdPost, err = client.CreatePostWithOptions(newPost, opts)
	} else {
		// Lexical形式またはコンテンツなしの場合は通常の作成
		createdPost, err = client.CreatePost(newPost)
	}

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created post: %s (ID: %s)", createdPost.Title, createdPost.ID))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(createdPost)
	}

	return nil
}

// PostsUpdateCmd is the command to update 投稿
type PostsUpdateCmd struct {
	ID       string `arg:"" help:"Post ID"`
	Title    string `help:"Post title" short:"t"`
	HTML     string `help:"Post content (HTML)" short:"c"`
	Markdown string `help:"Post content (Markdown)" short:"m"`
	Lexical  string `help:"Post content (Lexical JSON)" short:"x"`
	File     string `help:"Read content from file (auto-detect format)" type:"existingfile"`
	Status   string `help:"Post status (draft, published)"`
}

// Run executes the update subcommand of the posts command
func (c *PostsUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing post
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// コンテンツとフォーマットの決定
	var htmlContent string
	var format input.ContentFormat

	// ファイル指定の場合はフォーマット自動検出
	if c.File != "" {
		fileContent, detectedFormat, err := input.ReadContentWithFormat(c.File, "")
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		format = detectedFormat

		// フォーマットに応じて処理
		switch format {
		case input.FormatMarkdown:
			// Markdown→HTML変換
			htmlContent, err = markdown.ConvertToHTML(fileContent)
			if err != nil {
				return fmt.Errorf("failed to convert markdown to HTML: %w", err)
			}
		case input.FormatHTML:
			// HTMLはそのまま使用
			htmlContent = fileContent
		case input.FormatLexical:
			// Lexical JSONはそのまま使用（c.Lexicalに設定）
			c.Lexical = fileContent
		default:
			// 不明な形式の場合はHTMLとして扱う
			htmlContent = fileContent
		}
	} else {
		// インラインコンテンツの処理
		htmlContent = c.HTML

		// Markdownフラグが指定されている場合はMarkdown→HTML変換
		if c.Markdown != "" {
			htmlContent, err = markdown.ConvertToHTML(c.Markdown)
			if err != nil {
				return fmt.Errorf("failed to convert markdown to HTML: %w", err)
			}
		}
	}

	// Apply updates
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    existingPost.Status,
		UpdatedAt: existingPost.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	if c.Title != "" {
		updatePost.Title = c.Title
	}
	if htmlContent != "" {
		updatePost.HTML = htmlContent
	}
	if c.Lexical != "" {
		updatePost.Lexical = c.Lexical
	}
	if c.Status != "" {
		updatePost.Status = c.Status
	}

	// HTMLコンテンツが更新される場合は、自動的にsource=htmlを適用
	var updatedPost *ghostapi.Post
	if htmlContent != "" && c.Lexical == "" {
		// HTMLをサーバー側でLexical形式に変換
		opts := ghostapi.CreateOptions{
			Source: "html",
		}
		updatedPost, err = client.UpdatePostWithOptions(c.ID, updatePost, opts)
	} else {
		// Lexical形式またはHTMLの更新なしの場合は通常の更新
		updatedPost, err = client.UpdatePost(c.ID, updatePost)
	}

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated post: %s (ID: %s)", updatedPost.Title, updatedPost.ID))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(updatedPost)
	}

	return nil
}

// PostsDeleteCmd is the command to delete 投稿
type PostsDeleteCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run executes the delete subcommand of the posts command
func (c *PostsDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get post information to build confirmation message
	post, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete post '%s' (ID: %s)", post.Title, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Delete post
	if err := client.DeletePost(c.ID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted post (ID: %s)", c.ID))

	return nil
}

// PostsPublishCmd is the command to publish 下書き投稿
type PostsPublishCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run executes the publish subcommand of the posts command
func (c *PostsPublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing post
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Error if already published
	if existingPost.Status == "published" {
		return fmt.Errorf("post is already published")
	}

	// Change status to published
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    "published",
		UpdatedAt: existingPost.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update post
	publishedPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("failed to publish post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("published post: %s (ID: %s)", publishedPost.Title, publishedPost.ID))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(publishedPost)
	}

	return nil
}

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// PostsDraftsCmd is the command to retrieve draft 投稿 list
type PostsDraftsCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the drafts subcommand of the posts command
func (c *PostsDraftsCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get draft post list
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "draft",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("下書きfailed to list posts: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// Output in table format
	headers := []string{"ID", "Title", "Status", "Created", "Updated"}
	rows := make([][]string, len(response.Posts))
	for i, post := range response.Posts {
		rows[i] = []string{
			post.ID,
			post.Title,
			post.Status,
			post.CreatedAt.Format("2006-01-02"),
			post.UpdatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// PostsPublishedCmd is the command to retrieve published 投稿 list
type PostsPublishedCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the published subcommand of the posts command
func (c *PostsPublishedCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get published post list
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "published",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("公開済みfailed to list posts: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// Output in table format
	headers := []string{"ID", "Title", "Status", "Created", "Published"}
	rows := make([][]string, len(response.Posts))
	for i, post := range response.Posts {
		publishedAt := ""
		if post.PublishedAt != nil {
			publishedAt = post.PublishedAt.Format("2006-01-02")
		}
		rows[i] = []string{
			post.ID,
			post.Title,
			post.Status,
			post.CreatedAt.Format("2006-01-02"),
			publishedAt,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// PostsScheduledCmd is the command to retrieve scheduled 投稿 list
type PostsScheduledCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the scheduled subcommand of the posts command
func (c *PostsScheduledCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get scheduled post list
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "scheduled",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("予約failed to list posts: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// Output in table format
	headers := []string{"ID", "Title", "Status", "Created", "Scheduled"}
	rows := make([][]string, len(response.Posts))
	for i, post := range response.Posts {
		publishedAt := ""
		if post.PublishedAt != nil {
			publishedAt = post.PublishedAt.Format("2006-01-02 15:04")
		}
		rows[i] = []string{
			post.ID,
			post.Title,
			post.Status,
			post.CreatedAt.Format("2006-01-02"),
			publishedAt,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// ========================================
// Phase 1: URL取得
// ========================================

// PostsURLCmd is the command to get 投稿 web URL
type PostsURLCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
	Open     bool   `help:"Open URL in browser" short:"o"`
}

// Run executes the url subcommand of the posts command
func (c *PostsURLCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get post
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Get URL
	url := post.URL
	if url == "" {
		return fmt.Errorf("could not get post URL")
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

// PostsUnpublishCmd is the command to unpublish 公開済み投稿
type PostsUnpublishCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run executes the unpublish subcommand of the posts command
func (c *PostsUnpublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing post
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Error if already a draft
	if existingPost.Status == "draft" {
		return fmt.Errorf("post is already a draft")
	}

	// Change status to draft
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    "draft",
		UpdatedAt: existingPost.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update post
	unpublishedPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("failed to unpublish post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("unpublished post: %s (ID: %s)", unpublishedPost.Title, unpublishedPost.ID))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(unpublishedPost)
	}

	return nil
}

// ========================================
// Phase 3: 予約投稿
// ========================================

// PostsScheduleCmd is the command to schedule 投稿 for publishing
type PostsScheduleCmd struct {
	ID string `arg:"" help:"Post ID"`
	At string `help:"Schedule time (YYYY-MM-DD HH:MM)" required:""`
}

// Run executes the schedule subcommand of the posts command
func (c *PostsScheduleCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing post
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Parse datetime
	publishedAt, err := parseDateTime(c.At)
	if err != nil {
		return fmt.Errorf("failed to parse datetime: %w", err)
	}

	// Change status to scheduled and set publish date
	updatePost := &ghostapi.Post{
		Title:       existingPost.Title,
		Slug:        existingPost.Slug,
		HTML:        existingPost.HTML,
		Lexical:     existingPost.Lexical,
		Status:      "scheduled",
		PublishedAt: &publishedAt,
		UpdatedAt:   existingPost.UpdatedAt, // Use original updated_at from server (for optimistic locking)
	}

	// Update post
	scheduledPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("failed to schedule post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("scheduled post: %s (ID: %s, scheduled for: %s)",
			scheduledPost.Title, scheduledPost.ID, publishedAt.Format("2006-01-02 15:04")))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(scheduledPost)
	}

	return nil
}

// ========================================
// Phase 4: バッチ操作
// ========================================

// PostsBatchCmd はバッチ操作コマンドです
type PostsBatchCmd struct {
	Publish PostsBatchPublishCmd `cmd:"" help:"Batch publish posts"`
	Delete  PostsBatchDeleteCmd  `cmd:"" help:"Batch delete posts"`
}

// PostsBatchPublishCmd is the command to batch publish 投稿
type PostsBatchPublishCmd struct {
	IDs []string `arg:"" help:"Post IDs to publish"`
}

// Run executes the posts batch publish subcommand
func (c *PostsBatchPublishCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Publish each post
	successCount := 0
	for _, id := range c.IDs {
		// Get existing post
		existingPost, err := client.GetPost(id)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to get post (ID: %s): %v", id, err))
			continue
		}

		// すでに公開済みの場合はスキップ
		if existingPost.Status == "published" {
			formatter.PrintMessage(fmt.Sprintf("skipped (already published): %s (ID: %s)", existingPost.Title, id))
			continue
		}

		// Change status to published
		updatePost := &ghostapi.Post{
			Title:     existingPost.Title,
			Slug:      existingPost.Slug,
			HTML:      existingPost.HTML,
			Lexical:   existingPost.Lexical,
			Status:    "published",
			UpdatedAt: existingPost.UpdatedAt,
		}

		// Update post
		_, err = client.UpdatePost(id, updatePost)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to publish post (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("published: %s (ID: %s)", existingPost.Title, id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\ncompleted: %d件のpublished post", successCount))

	return nil
}

// PostsBatchDeleteCmd is the command to batch delete 投稿
type PostsBatchDeleteCmd struct {
	IDs []string `arg:"" help:"Post IDs to delete"`
}

// Run executes the posts batch delete subcommand
func (c *PostsBatchDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete %d posts", len(c.IDs))
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Delete each post
	successCount := 0
	for _, id := range c.IDs {
		// Delete post
		if err := client.DeletePost(id); err != nil {
			formatter.PrintMessage(fmt.Sprintf("failed to delete post (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("deleted (ID: %s)", id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\ncompleted: %d件のdeleted post", successCount))

	return nil
}

// ========================================
// Phase 4: 投稿検索
// ========================================

// PostsSearchCmd is the command to search 投稿
type PostsSearchCmd struct {
	Query string `arg:"" help:"Search query"`
	Limit int    `help:"Number of posts to retrieve" short:"l" aliases:"max,n" default:"15"`
}

// Run executes the search subcommand of the posts command
func (c *PostsSearchCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get post list（検索クエリはfilterとして渡す）
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "all",
		Limit:  c.Limit,
		Page:   1,
	})
	if err != nil {
		return fmt.Errorf("failed to search posts: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter posts matching query (simple implementation)
	var filteredPosts []ghostapi.Post
	for _, post := range response.Posts {
		if containsIgnoreCase(post.Title, c.Query) || containsIgnoreCase(post.HTML, c.Query) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(filteredPosts)
	}

	// Output in table format
	headers := []string{"ID", "Title", "Status", "Created"}
	rows := make([][]string, len(filteredPosts))
	for i, post := range filteredPosts {
		rows[i] = []string{
			post.ID,
			post.Title,
			post.Status,
			post.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// ========================================
// ヘルパー関数
// ========================================

// fileExists はファイルが存在するかチェックします
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// runCommand はコマンドを実行します
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// parseDateTime は日時文字列をパースします
func parseDateTime(s string) (time.Time, error) {
	// YYYY-MM-DD HH:MM形式をパース
	layout := "2006-01-02 15:04"
	return time.Parse(layout, s)
}

// containsIgnoreCase は大文字小文字を区別せずに部分文字列を検索します
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (strings.EqualFold(s, substr) || hasSubstringIgnoreCase(s, substr))
}

// hasSubstringIgnoreCase は大文字小文字を区別せずに部分文字列が含まれているかチェックします
func hasSubstringIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ========================================
// Phase 2: catコマンド
// ========================================

// PostsCatCmd is the command to show 投稿 content body
type PostsCatCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
	Format   string `help:"Output format (text, html, lexical)" default:"text"`
}

// Run executes the cat subcommand of the posts command
func (c *PostsCatCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get post
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output according to format
	var content string
	switch c.Format {
	case "html":
		content = post.HTML
	case "text":
		// Convert HTML to text
		content = html2text.HTML2Text(post.HTML)
	case "lexical":
		content = post.Lexical
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

// PostsCopyCmd is the command to copy 投稿
type PostsCopyCmd struct {
	IDOrSlug string `arg:"" help:"Source post ID or slug"`
	Title    string `help:"New title (defaults to 'Original Title (Copy)')" short:"t"`
}

// Run executes the copy subcommand of the posts command
func (c *PostsCopyCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get original post
	original, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Determine new title
	newTitle := c.Title
	if newTitle == "" {
		newTitle = original.Title + " (Copy)"
	}

	// Create new post (exclude ID/UUID/Slug/URL/dates, Status fixed to draft)
	newPost := &ghostapi.Post{
		Title:   newTitle,
		HTML:    original.HTML,
		Lexical: original.Lexical,
		Status:  "draft",
	}

	// Create post
	createdPost, err := client.CreatePost(newPost)
	if err != nil {
		return fmt.Errorf("failed to copy post: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("copied post: %s (ID: %s)", createdPost.Title, createdPost.ID))
	}

	// Also output post information if JSON format
	if root.JSON {
		return formatter.Print(createdPost)
	}

	return nil
}
