/**
 * users.go
 * ユーザー管理コマンド
 *
 * Ghostユーザー（サイト管理者・投稿者）の取得・更新機能を提供します。
 * ユーザーの作成・削除はGhostダッシュボードの招待機能を利用します。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// UsersCmd はユーザー管理コマンドです
type UsersCmd struct {
	List   UsersListCmd   `cmd:"" help:"List users"`
	Info   UsersInfoCmd   `cmd:"" help:"ユーザーの情報を表示"`
	Get    UsersInfoCmd   `cmd:"" hidden:"" help:"ユーザーを取得（非推奨: infoを使用してください）"`
	Update UsersUpdateCmd `cmd:"" help:"Update a user"`
}

// UsersListCmd はユーザー一覧を取得するコマンドです
type UsersListCmd struct {
	Limit   int    `help:"Number of users to retrieve" short:"l" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (e.g., roles,count.posts)" short:"i"`
	Filter  string `help:"Filter query"`
}

// Run はusersコマンドのlistサブコマンドを実行します
func (c *UsersListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ユーザー一覧を取得
	response, err := client.ListUsers(ghostapi.UserListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
		Filter:  c.Filter,
	})
	if err != nil {
		return fmt.Errorf("ユーザー一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Users)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Name", "Slug", "Email", "Created"}
	rows := make([][]string, len(response.Users))
	for i, user := range response.Users {
		rows[i] = []string{
			user.ID,
			user.Name,
			user.Slug,
			user.Email,
			user.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// UsersInfoCmd はユーザー情報を表示するコマンドです
type UsersInfoCmd struct {
	IDOrSlug string `arg:"" help:"User ID or slug (use 'slug:user-slug' format for slug)"`
}

// Run はusersコマンドのinfoサブコマンドを実行します
func (c *UsersInfoCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ユーザーを取得
	user, err := client.GetUser(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("ユーザーの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(user)
	}

	// キー/値形式で出力（ヘッダーなし）
	rows := [][]string{
		{"id", user.ID},
		{"name", user.Name},
		{"slug", user.Slug},
		{"email", user.Email},
		{"bio", user.Bio},
		{"location", user.Location},
		{"website", user.Website},
		{"profile_image", user.ProfileImage},
		{"cover_image", user.CoverImage},
		{"created", user.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", user.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	// ロール情報を追加
	if len(user.Roles) > 0 {
		roleNames := ""
		for i, role := range user.Roles {
			if i > 0 {
				roleNames += ", "
			}
			roleNames += role.Name
		}
		rows = append(rows, []string{"roles", roleNames})
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// UsersUpdateCmd はユーザーを更新するコマンドです
type UsersUpdateCmd struct {
	ID       string `arg:"" help:"User ID"`
	Name     string `help:"User name" short:"n"`
	Slug     string `help:"User slug"`
	Bio      string `help:"User bio" short:"b"`
	Location string `help:"User location" short:"l"`
	Website  string `help:"User website" short:"w"`
}

// Run はusersコマンドのupdateサブコマンドを実行します
func (c *UsersUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のユーザーを取得
	existingUser, err := client.GetUser(c.ID)
	if err != nil {
		return fmt.Errorf("ユーザーの取得に失敗: %w", err)
	}

	// 更新内容を反映
	updateUser := &ghostapi.User{
		Name:     existingUser.Name,
		Slug:     existingUser.Slug,
		Email:    existingUser.Email,
		Bio:      existingUser.Bio,
		Location: existingUser.Location,
		Website:  existingUser.Website,
	}

	if c.Name != "" {
		updateUser.Name = c.Name
	}
	if c.Slug != "" {
		updateUser.Slug = c.Slug
	}
	if c.Bio != "" {
		updateUser.Bio = c.Bio
	}
	if c.Location != "" {
		updateUser.Location = c.Location
	}
	if c.Website != "" {
		updateUser.Website = c.Website
	}

	// ユーザーを更新
	updatedUser, err := client.UpdateUser(c.ID, updateUser)
	if err != nil {
		return fmt.Errorf("ユーザーの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ユーザーを更新しました: %s (ID: %s)", updatedUser.Name, updatedUser.ID))
	}

	// JSON形式の場合はユーザー情報も出力
	if root.JSON {
		return formatter.Print(updatedUser)
	}

	return nil
}
