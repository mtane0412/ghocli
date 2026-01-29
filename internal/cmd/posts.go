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
	"time"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// PostsCmd は投稿管理コマンドです
type PostsCmd struct {
	List    PostsListCmd    `cmd:"" help:"List posts"`
	Get     PostsGetCmd     `cmd:"" help:"Get a post"`
	Create  PostsCreateCmd  `cmd:"" help:"Create a post"`
	Update  PostsUpdateCmd  `cmd:"" help:"Update a post"`
	Delete  PostsDeleteCmd  `cmd:"" help:"Delete a post"`
	Publish PostsPublishCmd `cmd:"" help:"Publish a draft"`
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

// PostsGetCmd は投稿を取得するコマンドです
type PostsGetCmd struct {
	IDOrSlug string `arg:"" help:"Post ID or slug"`
}

// Run はpostsコマンドのgetサブコマンドを実行します
func (c *PostsGetCmd) Run(root *RootFlags) error {
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
		UpdatedAt: time.Now(), // 更新時刻を設定
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

	// 確認なしで削除する場合を除き、確認を求める
	if !root.Force {
		// 投稿情報を取得して確認
		post, err := client.GetPost(c.ID)
		if err != nil {
			return fmt.Errorf("投稿の取得に失敗: %w", err)
		}

		fmt.Printf("本当に投稿「%s」(ID: %s)を削除しますか? [y/N]: ", post.Title, c.ID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("削除がキャンセルされました")
		}
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
		UpdatedAt: time.Now(),
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
