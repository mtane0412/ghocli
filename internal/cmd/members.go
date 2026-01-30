/**
 * members.go
 * メンバー管理コマンド
 *
 * Ghostメンバー（購読者）の作成、更新、削除機能を提供します。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// MembersCmd はメンバー管理コマンドです
type MembersCmd struct {
	List   MembersListCmd   `cmd:"" help:"List members"`
	Get    MembersGetCmd    `cmd:"" help:"Get a member"`
	Create MembersCreateCmd `cmd:"" help:"Create a member"`
	Update MembersUpdateCmd `cmd:"" help:"Update a member"`
	Delete MembersDeleteCmd `cmd:"" help:"Delete a member"`
}

// MembersListCmd はメンバー一覧を取得するコマンドです
type MembersListCmd struct {
	Limit  int    `help:"Number of members to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter query (e.g., status:paid)"`
	Order  string `help:"Sort order (e.g., created_at DESC)" short:"o"`
}

// Run はmembersコマンドのlistサブコマンドを実行します
func (c *MembersListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// メンバー一覧を取得
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
		Order:  c.Order,
	})
	if err != nil {
		return fmt.Errorf("メンバー一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Members)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Email", "Name", "Status", "Created"}
	rows := make([][]string, len(response.Members))
	for i, member := range response.Members {
		rows[i] = []string{
			member.ID,
			member.Email,
			member.Name,
			member.Status,
			member.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// MembersGetCmd はメンバーを取得するコマンドです
type MembersGetCmd struct {
	ID string `arg:"" help:"Member ID"`
}

// Run はmembersコマンドのgetサブコマンドを実行します
func (c *MembersGetCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// メンバーを取得
	member, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("メンバーの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(member)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", member.ID},
		{"UUID", member.UUID},
		{"Email", member.Email},
		{"Name", member.Name},
		{"Note", member.Note},
		{"Status", member.Status},
		{"Created", member.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", member.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	return formatter.PrintTable(headers, rows)
}

// MembersCreateCmd はメンバーを作成するコマンドです
type MembersCreateCmd struct {
	Email  string   `help:"Member email (required)" short:"e" required:""`
	Name   string   `help:"Member name" short:"n"`
	Note   string   `help:"Member note" short:"t"`
	Labels []string `help:"Member labels" short:"l"`
}

// Run はmembersコマンドのcreateサブコマンドを実行します
func (c *MembersCreateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 新規メンバーを作成
	newMember := &ghostapi.Member{
		Email: c.Email,
		Name:  c.Name,
		Note:  c.Note,
	}

	// ラベルを追加
	if len(c.Labels) > 0 {
		labels := make([]ghostapi.Label, len(c.Labels))
		for i, labelName := range c.Labels {
			labels[i] = ghostapi.Label{Name: labelName}
		}
		newMember.Labels = labels
	}

	createdMember, err := client.CreateMember(newMember)
	if err != nil {
		return fmt.Errorf("メンバーの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("メンバーを作成しました: %s (ID: %s)", createdMember.Email, createdMember.ID))
	}

	// JSON形式の場合はメンバー情報も出力
	if root.JSON {
		return formatter.Print(createdMember)
	}

	return nil
}

// MembersUpdateCmd はメンバーを更新するコマンドです
type MembersUpdateCmd struct {
	ID     string   `arg:"" help:"Member ID"`
	Name   string   `help:"Member name" short:"n"`
	Note   string   `help:"Member note" short:"t"`
	Labels []string `help:"Member labels" short:"l"`
}

// Run はmembersコマンドのupdateサブコマンドを実行します
func (c *MembersUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のメンバーを取得
	existingMember, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("メンバーの取得に失敗: %w", err)
	}

	// 更新内容を反映
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: existingMember.Labels,
	}

	if c.Name != "" {
		updateMember.Name = c.Name
	}
	if c.Note != "" {
		updateMember.Note = c.Note
	}
	if len(c.Labels) > 0 {
		labels := make([]ghostapi.Label, len(c.Labels))
		for i, labelName := range c.Labels {
			labels[i] = ghostapi.Label{Name: labelName}
		}
		updateMember.Labels = labels
	}

	// メンバーを更新
	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("メンバーの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("メンバーを更新しました: %s (ID: %s)", updatedMember.Email, updatedMember.ID))
	}

	// JSON形式の場合はメンバー情報も出力
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersDeleteCmd はメンバーを削除するコマンドです
type MembersDeleteCmd struct {
	ID string `arg:"" help:"Member ID"`
}

// Run はmembersコマンドのdeleteサブコマンドを実行します
func (c *MembersDeleteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// メンバー情報を取得して確認メッセージを構築
	member, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("メンバーの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("delete member '%s' (ID: %s)", member.Email, c.ID)
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
	}

	// メンバーを削除
	if err := client.DeleteMember(c.ID); err != nil {
		return fmt.Errorf("メンバーの削除に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	formatter.PrintMessage(fmt.Sprintf("メンバーを削除しました (ID: %s)", c.ID))

	return nil
}
