/**
 * tags.go
 * タグ管理コマンド
 *
 * Ghostタグの作成、更新、削除機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/fields"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// TagsCmd はタグ管理コマンドです
type TagsCmd struct {
	List   TagsListCmd   `cmd:"" help:"List tags"`
	Get    TagsInfoCmd   `cmd:"" help:"タグの情報を表示"`
	Create TagsCreateCmd `cmd:"" help:"Create a tag"`
	Update TagsUpdateCmd `cmd:"" help:"Update a tag"`
	Delete TagsDeleteCmd `cmd:"" help:"Delete a tag"`
}

// TagsListCmd はタグ一覧を取得するコマンドです
type TagsListCmd struct {
	Limit   int    `help:"Number of tags to retrieve" short:"l" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (count.posts)" short:"i"`
}

// Run はtagsコマンドのlistサブコマンドを実行します
func (c *TagsListCmd) Run(ctx context.Context, root *RootFlags) error {
	// JSON単独（--fieldsなし）の場合は利用可能なフィールド一覧を表示
	if root.JSON && root.Fields == "" {
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fields.ListAvailable(fields.TagFields))
		return nil
	}

	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.TagFields)
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

	// タグ一覧を取得
	response, err := client.ListTags(ghostapi.TagListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
	})
	if err != nil {
		return fmt.Errorf("タグ一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Tag構造体をmap[string]interface{}に変換
		var tagsData []map[string]interface{}
		for _, tag := range response.Tags {
			tagMap, err := outfmt.StructToMap(tag)
			if err != nil {
				return fmt.Errorf("タグデータの変換に失敗: %w", err)
			}
			tagsData = append(tagsData, tagMap)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, tagsData, selectedFields)
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Tags)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Name", "Slug", "Visibility", "Created"}
	rows := make([][]string, len(response.Tags))
	for i, tag := range response.Tags {
		rows[i] = []string{
			tag.ID,
			tag.Name,
			tag.Slug,
			tag.Visibility,
			tag.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// TagsInfoCmd はタグ情報を表示するコマンドです
type TagsInfoCmd struct {
	IDOrSlug string `arg:"" help:"Tag ID or slug (use 'slug:tag-name' format for slug)"`
}

// Run はtagsコマンドのinfoサブコマンドを実行します
func (c *TagsInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// JSON単独（--fieldsなし）の場合は利用可能なフィールド一覧を表示
	if root.JSON && root.Fields == "" {
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fields.ListAvailable(fields.TagFields))
		return nil
	}

	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.TagFields)
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

	// タグを取得
	tag, err := client.GetTag(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("タグの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Tag構造体をmap[string]interface{}に変換
		tagMap, err := outfmt.StructToMap(tag)
		if err != nil {
			return fmt.Errorf("タグデータの変換に失敗: %w", err)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, []map[string]interface{}{tagMap}, selectedFields)
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(tag)
	}

	// キー/値形式で出力（ヘッダーなし）
	rows := [][]string{
		{"id", tag.ID},
		{"name", tag.Name},
		{"slug", tag.Slug},
		{"description", tag.Description},
		{"visibility", tag.Visibility},
		{"created", tag.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", tag.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// TagsCreateCmd はタグを作成するコマンドです
type TagsCreateCmd struct {
	Name        string `help:"Tag name" short:"n" required:""`
	Description string `help:"Tag description" short:"d"`
	Visibility  string `help:"Tag visibility (public, internal)" default:"public"`
}

// Run はtagsコマンドのcreateサブコマンドを実行します
func (c *TagsCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 新規タグを作成
	newTag := &ghostapi.Tag{
		Name:        c.Name,
		Description: c.Description,
		Visibility:  c.Visibility,
	}

	createdTag, err := client.CreateTag(newTag)
	if err != nil {
		return fmt.Errorf("タグの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("タグを作成しました: %s (ID: %s)", createdTag.Name, createdTag.ID))
	}

	// JSON形式の場合はタグ情報も出力
	if root.JSON {
		return formatter.Print(createdTag)
	}

	return nil
}

// TagsUpdateCmd はタグを更新するコマンドです
type TagsUpdateCmd struct {
	ID          string `arg:"" help:"Tag ID"`
	Name        string `help:"Tag name" short:"n"`
	Description string `help:"Tag description" short:"d"`
	Visibility  string `help:"Tag visibility (public, internal)"`
}

// Run はtagsコマンドのupdateサブコマンドを実行します
func (c *TagsUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のタグを取得
	existingTag, err := client.GetTag(c.ID)
	if err != nil {
		return fmt.Errorf("タグの取得に失敗: %w", err)
	}

	// 更新内容を反映
	updateTag := &ghostapi.Tag{
		Name:        existingTag.Name,
		Slug:        existingTag.Slug,
		Description: existingTag.Description,
		Visibility:  existingTag.Visibility,
	}

	if c.Name != "" {
		updateTag.Name = c.Name
	}
	if c.Description != "" {
		updateTag.Description = c.Description
	}
	if c.Visibility != "" {
		updateTag.Visibility = c.Visibility
	}

	// タグを更新
	updatedTag, err := client.UpdateTag(c.ID, updateTag)
	if err != nil {
		return fmt.Errorf("タグの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("タグを更新しました: %s (ID: %s)", updatedTag.Name, updatedTag.ID))
	}

	// JSON形式の場合はタグ情報も出力
	if root.JSON {
		return formatter.Print(updatedTag)
	}

	return nil
}

// TagsDeleteCmd はタグを削除するコマンドです
type TagsDeleteCmd struct {
	ID string `arg:"" help:"Tag ID"`
}

// Run はtagsコマンドのdeleteサブコマンドを実行します
func (c *TagsDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// タグ情報を取得して確認メッセージを構築
	tag, err := client.GetTag(c.ID)
	if err != nil {
		return fmt.Errorf("タグの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete tag '%s' (ID: %s)", tag.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// タグを削除
	if err := client.DeleteTag(c.ID); err != nil {
		return fmt.Errorf("タグの削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("タグを削除しました (ID: %s)", c.ID))

	return nil
}
