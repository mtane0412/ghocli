/**
 * posts.go
 * 投稿管理コマンド
 *
 * Ghost投稿の作成、更新、削除、公開機能を提供します。
 */

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// PostsCmd は投稿管理コマンドです
type PostsCmd struct {
	List    PostsListCmd    `cmd:"" help:"List posts"`
	Info    PostsInfoCmd    `cmd:"" help:"投稿の情報を表示"`
	Get     PostsInfoCmd    `cmd:"" hidden:"" help:"投稿を取得（非推奨: infoを使用してください）"`
	Cat     PostsCatCmd     `cmd:"" help:"本文コンテンツを表示"`
	Create  PostsCreateCmd  `cmd:"" help:"Create a post"`
	Update  PostsUpdateCmd  `cmd:"" help:"Update a post"`
	Delete  PostsDeleteCmd  `cmd:"" help:"Delete a post"`
	Publish PostsPublishCmd `cmd:"" help:"Publish a draft"`

	// Phase 1: ステータス別一覧ショートカット
	Drafts    PostsDraftsCmd    `cmd:"" help:"List draft posts"`
	Published PostsPublishedCmd `cmd:"" help:"List published posts"`
	Scheduled PostsScheduledCmd `cmd:"" help:"List scheduled posts"`

	// Phase 1: URL取得
	URL PostsURLCmd `cmd:"" help:"Get post URL"`

	// Phase 2: 状態変更
	Unpublish PostsUnpublishCmd `cmd:"" help:"Unpublish a post"`

	// Phase 3: 予約投稿
	Schedule PostsScheduleCmd `cmd:"" help:"Schedule a post"`

	// Phase 4: バッチ操作
	Batch  PostsBatchCmd  `cmd:"" help:"Batch operations"`
	Search PostsSearchCmd `cmd:"" help:"Search posts"`

	// Phase 8.3: コピー
	Copy PostsCopyCmd `cmd:"" help:"投稿をコピー"`
}

// PostsListCmd は投稿一覧を取得するコマンドです
type PostsListCmd struct {
	Status string `help:"Filter by status (draft, published, scheduled, all)" short:"S" default:"all"`
	Limit  int    `help:"Number of posts to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
}

// Run はpostsコマンドのlistサブコマンドを実行します
func (c *PostsListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿一覧を取得
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: c.Status,
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("投稿一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// テーブル形式で出力
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

// PostsInfoCmd は投稿情報を表示するコマンドです
type PostsInfoCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
}

// Run はpostsコマンドのinfoサブコマンドを実行します
func (c *PostsInfoCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿を取得
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(post)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", post.ID},
		{"Title", post.Title},
		{"Slug", post.Slug},
		{"Status", post.Status},
		{"Created", post.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", post.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if post.PublishedAt != nil {
		rows = append(rows, []string{"Published", post.PublishedAt.Format("2006-01-02 15:04:05")})
	}

	return formatter.PrintTable(headers, rows)
}

// PostsCreateCmd は投稿を作成するコマンドです
type PostsCreateCmd struct {
	Title   string `help:"Post title" short:"t" required:""`
	HTML    string `help:"Post content (HTML)" short:"c"`
	Lexical string `help:"Post content (Lexical JSON)" short:"x"`
	Status  string `help:"Post status (draft, published)" default:"draft"`
}

// Run はpostsコマンドのcreateサブコマンドを実行します
func (c *PostsCreateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 新規投稿を作成
	newPost := &ghostapi.Post{
		Title:   c.Title,
		HTML:    c.HTML,
		Lexical: c.Lexical,
		Status:  c.Status,
	}

	createdPost, err := client.CreatePost(newPost)
	if err != nil {
		return fmt.Errorf("投稿の作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿を作成しました: %s (ID: %s)", createdPost.Title, createdPost.ID))
	}

	// JSON形式の場合は投稿情報も出力
	if root.JSON {
		return formatter.Print(createdPost)
	}

	return nil
}

// PostsUpdateCmd は投稿を更新するコマンドです
type PostsUpdateCmd struct {
	ID      string `arg:"" help:"Post ID"`
	Title   string `help:"Post title" short:"t"`
	HTML    string `help:"Post content (HTML)" short:"c"`
	Lexical string `help:"Post content (Lexical JSON)" short:"x"`
	Status  string `help:"Post status (draft, published)"`
}

// Run はpostsコマンドのupdateサブコマンドを実行します
func (c *PostsUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存の投稿を取得
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 更新内容を反映
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    existingPost.Status,
		UpdatedAt: existingPost.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	if c.Title != "" {
		updatePost.Title = c.Title
	}
	if c.HTML != "" {
		updatePost.HTML = c.HTML
	}
	if c.Lexical != "" {
		updatePost.Lexical = c.Lexical
	}
	if c.Status != "" {
		updatePost.Status = c.Status
	}

	// 投稿を更新
	updatedPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("投稿の更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿を更新しました: %s (ID: %s)", updatedPost.Title, updatedPost.ID))
	}

	// JSON形式の場合は投稿情報も出力
	if root.JSON {
		return formatter.Print(updatedPost)
	}

	return nil
}

// PostsDeleteCmd は投稿を削除するコマンドです
type PostsDeleteCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run はpostsコマンドのdeleteサブコマンドを実行します
func (c *PostsDeleteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿情報を取得して確認メッセージを構築
	post, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete post '%s' (ID: %s)", post.Title, c.ID)
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
	}

	// 投稿を削除
	if err := client.DeletePost(c.ID); err != nil {
		return fmt.Errorf("投稿の削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("投稿を削除しました (ID: %s)", c.ID))

	return nil
}

// PostsPublishCmd は下書き投稿を公開するコマンドです
type PostsPublishCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run はpostsコマンドのpublishサブコマンドを実行します
func (c *PostsPublishCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存の投稿を取得
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// すでに公開済みの場合はエラー
	if existingPost.Status == "published" {
		return fmt.Errorf("この投稿はすでに公開されています")
	}

	// ステータスをpublishedに変更
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    "published",
		UpdatedAt: existingPost.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// 投稿を更新
	publishedPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("投稿の公開に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿を公開しました: %s (ID: %s)", publishedPost.Title, publishedPost.ID))
	}

	// JSON形式の場合は投稿情報も出力
	if root.JSON {
		return formatter.Print(publishedPost)
	}

	return nil
}

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// PostsDraftsCmd は下書き投稿一覧を取得するコマンドです
type PostsDraftsCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpostsコマンドのdraftsサブコマンドを実行します
func (c *PostsDraftsCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 下書き投稿一覧を取得
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "draft",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("下書き投稿一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// テーブル形式で出力
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

// PostsPublishedCmd は公開済み投稿一覧を取得するコマンドです
type PostsPublishedCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpostsコマンドのpublishedサブコマンドを実行します
func (c *PostsPublishedCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 公開済み投稿一覧を取得
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "published",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("公開済み投稿一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// テーブル形式で出力
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

// PostsScheduledCmd は予約投稿一覧を取得するコマンドです
type PostsScheduledCmd struct {
	Limit int `help:"Number of posts to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はpostsコマンドのscheduledサブコマンドを実行します
func (c *PostsScheduledCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 予約投稿一覧を取得
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "scheduled",
		Limit:  c.Limit,
		Page:   c.Page,
	})
	if err != nil {
		return fmt.Errorf("予約投稿一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Posts)
	}

	// テーブル形式で出力
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

// PostsURLCmd は投稿のWeb URLを取得するコマンドです
type PostsURLCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
	Open     bool   `help:"Open URL in browser" short:"o"`
}

// Run はpostsコマンドのurlサブコマンドを実行します
func (c *PostsURLCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿を取得
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// URLを取得
	url := post.URL
	if url == "" {
		return fmt.Errorf("投稿のURLが取得できませんでした")
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

// PostsUnpublishCmd は公開済み投稿を下書きに戻すコマンドです
type PostsUnpublishCmd struct {
	ID string `arg:"" help:"Post ID"`
}

// Run はpostsコマンドのunpublishサブコマンドを実行します
func (c *PostsUnpublishCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存の投稿を取得
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// すでに下書きの場合はエラー
	if existingPost.Status == "draft" {
		return fmt.Errorf("この投稿はすでに下書きです")
	}

	// ステータスをdraftに変更
	updatePost := &ghostapi.Post{
		Title:     existingPost.Title,
		Slug:      existingPost.Slug,
		HTML:      existingPost.HTML,
		Lexical:   existingPost.Lexical,
		Status:    "draft",
		UpdatedAt: existingPost.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// 投稿を更新
	unpublishedPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("投稿の非公開化に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿を下書きに戻しました: %s (ID: %s)", unpublishedPost.Title, unpublishedPost.ID))
	}

	// JSON形式の場合は投稿情報も出力
	if root.JSON {
		return formatter.Print(unpublishedPost)
	}

	return nil
}

// ========================================
// Phase 3: 予約投稿
// ========================================

// PostsScheduleCmd は投稿を予約公開に設定するコマンドです
type PostsScheduleCmd struct {
	ID string `arg:"" help:"Post ID"`
	At string `help:"Schedule time (YYYY-MM-DD HH:MM)" required:""`
}

// Run はpostsコマンドのscheduleサブコマンドを実行します
func (c *PostsScheduleCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存の投稿を取得
	existingPost, err := client.GetPost(c.ID)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 日時をパース
	publishedAt, err := parseDateTime(c.At)
	if err != nil {
		return fmt.Errorf("日時のパースに失敗: %w", err)
	}

	// ステータスをscheduledに変更し、公開日時を設定
	updatePost := &ghostapi.Post{
		Title:       existingPost.Title,
		Slug:        existingPost.Slug,
		HTML:        existingPost.HTML,
		Lexical:     existingPost.Lexical,
		Status:      "scheduled",
		PublishedAt: &publishedAt,
		UpdatedAt:   existingPost.UpdatedAt, // サーバーから取得した元のupdated_atを使用（楽観的ロックのため）
	}

	// 投稿を更新
	scheduledPost, err := client.UpdatePost(c.ID, updatePost)
	if err != nil {
		return fmt.Errorf("投稿の予約公開設定に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿を予約公開に設定しました: %s (ID: %s, 公開予定: %s)",
			scheduledPost.Title, scheduledPost.ID, publishedAt.Format("2006-01-02 15:04")))
	}

	// JSON形式の場合は投稿情報も出力
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

// PostsBatchPublishCmd は複数投稿を一括公開するコマンドです
type PostsBatchPublishCmd struct {
	IDs []string `arg:"" help:"Post IDs to publish"`
}

// Run はposts batch publishサブコマンドを実行します
func (c *PostsBatchPublishCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 各投稿を公開
	successCount := 0
	for _, id := range c.IDs {
		// 既存の投稿を取得
		existingPost, err := client.GetPost(id)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("投稿の取得に失敗 (ID: %s): %v", id, err))
			continue
		}

		// すでに公開済みの場合はスキップ
		if existingPost.Status == "published" {
			formatter.PrintMessage(fmt.Sprintf("スキップ (すでに公開済み): %s (ID: %s)", existingPost.Title, id))
			continue
		}

		// ステータスをpublishedに変更
		updatePost := &ghostapi.Post{
			Title:     existingPost.Title,
			Slug:      existingPost.Slug,
			HTML:      existingPost.HTML,
			Lexical:   existingPost.Lexical,
			Status:    "published",
			UpdatedAt: existingPost.UpdatedAt,
		}

		// 投稿を更新
		_, err = client.UpdatePost(id, updatePost)
		if err != nil {
			formatter.PrintMessage(fmt.Sprintf("投稿の公開に失敗 (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("公開しました: %s (ID: %s)", existingPost.Title, id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\n完了: %d件の投稿を公開しました", successCount))

	return nil
}

// PostsBatchDeleteCmd は複数投稿を一括削除するコマンドです
type PostsBatchDeleteCmd struct {
	IDs []string `arg:"" help:"Post IDs to delete"`
}

// Run はposts batch deleteサブコマンドを実行します
func (c *PostsBatchDeleteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete %d posts", len(c.IDs))
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 各投稿を削除
	successCount := 0
	for _, id := range c.IDs {
		// 投稿を削除
		if err := client.DeletePost(id); err != nil {
			formatter.PrintMessage(fmt.Sprintf("投稿の削除に失敗 (ID: %s): %v", id, err))
			continue
		}

		formatter.PrintMessage(fmt.Sprintf("削除しました (ID: %s)", id))
		successCount++
	}

	formatter.PrintMessage(fmt.Sprintf("\n完了: %d件の投稿を削除しました", successCount))

	return nil
}

// ========================================
// Phase 4: 投稿検索
// ========================================

// PostsSearchCmd は投稿を検索するコマンドです
type PostsSearchCmd struct {
	Query string `arg:"" help:"Search query"`
	Limit int    `help:"Number of posts to retrieve" short:"l" default:"15"`
}

// Run はpostsコマンドのsearchサブコマンドを実行します
func (c *PostsSearchCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿一覧を取得（検索クエリはfilterとして渡す）
	response, err := client.ListPosts(ghostapi.ListOptions{
		Status: "all",
		Limit:  c.Limit,
		Page:   1,
	})
	if err != nil {
		return fmt.Errorf("投稿検索に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// クエリに一致する投稿をフィルタリング（簡易的な実装）
	var filteredPosts []ghostapi.Post
	for _, post := range response.Posts {
		if containsIgnoreCase(post.Title, c.Query) || containsIgnoreCase(post.HTML, c.Query) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(filteredPosts)
	}

	// テーブル形式で出力
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

// PostsCatCmd は投稿の本文コンテンツを表示するコマンドです
type PostsCatCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
	Format   string `help:"Output format (html, text, lexical)" default:"html"`
}

// Run はpostsコマンドのcatサブコマンドを実行します
func (c *PostsCatCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 投稿を取得
	post, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// フォーマットに応じて出力
	var content string
	switch c.Format {
	case "html":
		content = post.HTML
	case "text":
		// TODO: HTMLからテキストへの変換を実装
		// 現時点ではHTMLをそのまま出力
		content = post.HTML
	case "lexical":
		content = post.Lexical
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

// PostsCopyCmd は投稿をコピーするコマンドです
type PostsCopyCmd struct {
	IDOrSlug string `arg:"" help:"コピー元の投稿ID またはスラッグ"`
	Title    string `help:"新しいタイトル（省略時は '元タイトル (Copy)'）" short:"t"`
}

// Run はpostsコマンドのcopyサブコマンドを実行します
func (c *PostsCopyCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 元の投稿を取得
	original, err := client.GetPost(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	// 新しいタイトルを決定
	newTitle := c.Title
	if newTitle == "" {
		newTitle = original.Title + " (Copy)"
	}

	// 新しい投稿を作成（ID/UUID/Slug/URL/日時は除外、Statusはdraft固定）
	newPost := &ghostapi.Post{
		Title:   newTitle,
		HTML:    original.HTML,
		Lexical: original.Lexical,
		Status:  "draft",
	}

	// 投稿を作成
	createdPost, err := client.CreatePost(newPost)
	if err != nil {
		return fmt.Errorf("投稿のコピーに失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("投稿をコピーしました: %s (ID: %s)", createdPost.Title, createdPost.ID))
	}

	// JSON形式の場合は投稿情報も出力
	if root.JSON {
		return formatter.Print(createdPost)
	}

	return nil
}
