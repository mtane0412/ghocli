/**
 * members.go
 * メンバー管理コマンド
 *
 * Ghostメンバー（購読者）の作成、更新、削除機能を提供します。
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

// MembersCmd はメンバー管理コマンドです
type MembersCmd struct {
	List   MembersListCmd   `cmd:"" help:"List members"`
	Get    MembersInfoCmd   `cmd:"" help:"メンバーの情報を表示"`
	Create MembersCreateCmd `cmd:"" help:"Create a member"`
	Update MembersUpdateCmd `cmd:"" help:"Update a member"`
	Delete MembersDeleteCmd `cmd:"" help:"Delete a member"`

	// Phase 1: ステータス別一覧ショートカット
	Paid MembersPaidCmd `cmd:"" help:"List paid members"`
	Free MembersFreeCmd `cmd:"" help:"List free members"`

	// Phase 3: ラベル操作
	Label   MembersLabelCmd   `cmd:"" help:"Add label to member"`
	Unlabel MembersUnlabelCmd `cmd:"" help:"Remove label from member"`
	Recent  MembersRecentCmd  `cmd:"" help:"List recently created members"`
}

// MembersListCmd はメンバー一覧を取得するコマンドです
type MembersListCmd struct {
	Limit  int    `help:"Number of members to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter query (e.g., status:paid)"`
	Order  string `help:"Sort order (e.g., created_at DESC)" short:"o"`
}

// Run はmembersコマンドのlistサブコマンドを実行します
func (c *MembersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// JSON単独（--fieldsなし）の場合は利用可能なフィールド一覧を表示
	if root.JSON && root.Fields == "" {
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fields.ListAvailable(fields.MemberFields))
		return nil
	}

	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.MemberFields)
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

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Member構造体をmap[string]interface{}に変換
		var membersData []map[string]interface{}
		for _, member := range response.Members {
			memberMap, err := outfmt.StructToMap(member)
			if err != nil {
				return fmt.Errorf("メンバーデータの変換に失敗: %w", err)
			}
			membersData = append(membersData, memberMap)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, membersData, selectedFields)
	}

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

// MembersInfoCmd はメンバー情報を表示するコマンドです
type MembersInfoCmd struct {
	ID string `arg:"" help:"Member ID"`
}

// Run はmembersコマンドのinfoサブコマンドを実行します
func (c *MembersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// JSON単独（--fieldsなし）の場合は利用可能なフィールド一覧を表示
	if root.JSON && root.Fields == "" {
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fields.ListAvailable(fields.MemberFields))
		return nil
	}

	// フィールド指定をパース
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.MemberFields)
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

	// メンバーを取得
	member, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("メンバーの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// フィールド指定がある場合はフィルタリングして出力
	if len(selectedFields) > 0 {
		// Member構造体をmap[string]interface{}に変換
		memberMap, err := outfmt.StructToMap(member)
		if err != nil {
			return fmt.Errorf("メンバーデータの変換に失敗: %w", err)
		}

		// フィールドフィルタリングして出力
		return outfmt.FilterFields(formatter, []map[string]interface{}{memberMap}, selectedFields)
	}

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(member)
	}

	// キー/値形式で出力（ヘッダーなし）
	rows := [][]string{
		{"id", member.ID},
		{"uuid", member.UUID},
		{"email", member.Email},
		{"name", member.Name},
		{"note", member.Note},
		{"status", member.Status},
		{"created", member.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", member.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// MembersCreateCmd はメンバーを作成するコマンドです
type MembersCreateCmd struct {
	Email  string   `help:"Member email (required)" short:"e" required:""`
	Name   string   `help:"Member name" short:"n"`
	Note   string   `help:"Member note" short:"t"`
	Labels []string `help:"Member labels" short:"l"`
}

// Run はmembersコマンドのcreateサブコマンドを実行します
func (c *MembersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *MembersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
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
func (c *MembersDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
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
	if err := ConfirmDestructive(ctx, root, action); err != nil {
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

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// MembersPaidCmd は有料会員一覧を取得するコマンドです
type MembersPaidCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はmembersコマンドのpaidサブコマンドを実行します
func (c *MembersPaidCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 有料会員一覧を取得
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: "status:paid",
	})
	if err != nil {
		return fmt.Errorf("有料会員一覧の取得に失敗: %w", err)
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

// MembersFreeCmd は無料会員一覧を取得するコマンドです
type MembersFreeCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run はmembersコマンドのfreeサブコマンドを実行します
func (c *MembersFreeCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 無料会員一覧を取得
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: "status:free",
	})
	if err != nil {
		return fmt.Errorf("無料会員一覧の取得に失敗: %w", err)
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

// ========================================
// Phase 3: ラベル操作
// ========================================

// MembersLabelCmd はメンバーにラベルを追加するコマンドです
type MembersLabelCmd struct {
	ID    string `arg:"" help:"Member ID"`
	Label string `arg:"" help:"Label name"`
}

// Run はmembersコマンドのlabelサブコマンドを実行します
func (c *MembersLabelCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// 既存のラベルにLabel名がある場合はスキップ
	for _, label := range existingMember.Labels {
		if label.Name == c.Label {
			// 出力フォーマッターを作成
			formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
			formatter.PrintMessage(fmt.Sprintf("メンバーはすでにラベル '%s' を持っています (ID: %s)", c.Label, c.ID))
			return nil
		}
	}

	// 既存のラベルに新しいラベルを追加
	newLabels := append(existingMember.Labels, ghostapi.Label{Name: c.Label})

	// メンバーを更新
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: newLabels,
	}

	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("メンバーの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("メンバーにラベルを追加しました: %s (ID: %s, Label: %s)", updatedMember.Email, updatedMember.ID, c.Label))
	}

	// JSON形式の場合はメンバー情報も出力
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersUnlabelCmd はメンバーからラベルを削除するコマンドです
type MembersUnlabelCmd struct {
	ID    string `arg:"" help:"Member ID"`
	Label string `arg:"" help:"Label name"`
}

// Run はmembersコマンドのunlabelサブコマンドを実行します
func (c *MembersUnlabelCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// 既存のラベルから指定されたラベルを削除
	var newLabels []ghostapi.Label
	found := false
	for _, label := range existingMember.Labels {
		if label.Name != c.Label {
			newLabels = append(newLabels, label)
		} else {
			found = true
		}
	}

	// ラベルが見つからなかった場合
	if !found {
		// 出力フォーマッターを作成
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fmt.Sprintf("メンバーはラベル '%s' を持っていません (ID: %s)", c.Label, c.ID))
		return nil
	}

	// メンバーを更新
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: newLabels,
	}

	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("メンバーの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("メンバーからラベルを削除しました: %s (ID: %s, Label: %s)", updatedMember.Email, updatedMember.ID, c.Label))
	}

	// JSON形式の場合はメンバー情報も出力
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersRecentCmd は最近登録したメンバー一覧を取得するコマンドです
type MembersRecentCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" default:"15"`
}

// Run はmembersコマンドのrecentサブコマンドを実行します
func (c *MembersRecentCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 最近登録したメンバー一覧を取得（created_atの降順でソート）
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit: c.Limit,
		Page:  1,
		Order: "created_at DESC",
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
			member.CreatedAt.Format("2006-01-02 15:04"),
		}
	}

	return formatter.PrintTable(headers, rows)
}
