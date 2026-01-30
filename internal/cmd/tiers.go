/**
 * tiers.go
 * ティア管理コマンド
 *
 * Ghostティアの管理機能を提供します。
 * Create/Update操作には確認機構が適用されます。
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// TiersCmd はティア管理コマンドです
type TiersCmd struct {
	List   TiersListCmd   `cmd:"" help:"List tiers"`
	Get    TiersGetCmd    `cmd:"" help:"Get a tier"`
	Create TiersCreateCmd `cmd:"" help:"Create a tier"`
	Update TiersUpdateCmd `cmd:"" help:"Update a tier"`
}

// TiersListCmd はティア一覧を取得するコマンドです
type TiersListCmd struct {
	Limit   int    `help:"Number of tiers to retrieve" short:"l" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (monthly_price,yearly_price,benefits)" short:"i"`
	Filter  string `help:"Filter condition"`
}

// Run はtiersコマンドのlistサブコマンドを実行します
func (c *TiersListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ティア一覧を取得
	response, err := client.ListTiers(ghostapi.TierListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
		Filter:  c.Filter,
	})
	if err != nil {
		return fmt.Errorf("ティア一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Tiers)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Name", "Slug", "Type", "Active", "Visibility", "Created"}
	rows := make([][]string, len(response.Tiers))
	for i, tier := range response.Tiers {
		active := "false"
		if tier.Active {
			active = "true"
		}
		rows[i] = []string{
			tier.ID,
			tier.Name,
			tier.Slug,
			tier.Type,
			active,
			tier.Visibility,
			tier.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// TiersGetCmd はティアを取得するコマンドです
type TiersGetCmd struct {
	IDOrSlug string `arg:"" help:"Tier ID or slug (use 'slug:tier-name' format for slug)"`
	Include  string `help:"Include additional data (monthly_price,yearly_price,benefits)" short:"i"`
}

// Run はtiersコマンドのgetサブコマンドを実行します
func (c *TiersGetCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// ティアを取得
	tier, err := client.GetTier(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("ティアの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(tier)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", tier.ID},
		{"Name", tier.Name},
		{"Slug", tier.Slug},
		{"Description", tier.Description},
		{"Type", tier.Type},
		{"Active", fmt.Sprintf("%t", tier.Active)},
		{"Visibility", tier.Visibility},
		{"Welcome Page URL", tier.WelcomePageURL},
		{"Monthly Price", fmt.Sprintf("%d", tier.MonthlyPrice)},
		{"Yearly Price", fmt.Sprintf("%d", tier.YearlyPrice)},
		{"Currency", tier.Currency},
		{"Benefits", strings.Join(tier.Benefits, ", ")},
		{"Created", tier.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", tier.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	return formatter.PrintTable(headers, rows)
}

// TiersCreateCmd はティアを作成するコマンドです
type TiersCreateCmd struct {
	Name           string   `help:"Tier name" short:"n" required:""`
	Description    string   `help:"Tier description" short:"d"`
	Type           string   `help:"Tier type (free, paid)" default:"paid"`
	Visibility     string   `help:"Visibility (public, none)" default:"public"`
	MonthlyPrice   int      `help:"Monthly price (in smallest currency unit)"`
	YearlyPrice    int      `help:"Yearly price (in smallest currency unit)"`
	Currency       string   `help:"Currency code (e.g., JPY, USD)" default:"JPY"`
	WelcomePageURL string   `help:"Welcome page URL"`
	Benefits       []string `help:"Benefits list" short:"b"`
}

// Run はtiersコマンドのcreateサブコマンドを実行します
func (c *TiersCreateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 破壊的操作の確認
	priceInfo := ""
	if c.Type == "paid" {
		priceInfo = fmt.Sprintf(" (monthly: %d %s, yearly: %d %s)", c.MonthlyPrice, c.Currency, c.YearlyPrice, c.Currency)
	}
	action := fmt.Sprintf("create tier '%s'%s", c.Name, priceInfo)
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
	}

	// 新規ティアを作成
	newTier := &ghostapi.Tier{
		Name:           c.Name,
		Description:    c.Description,
		Type:           c.Type,
		Visibility:     c.Visibility,
		MonthlyPrice:   c.MonthlyPrice,
		YearlyPrice:    c.YearlyPrice,
		Currency:       c.Currency,
		WelcomePageURL: c.WelcomePageURL,
		Benefits:       c.Benefits,
	}

	createdTier, err := client.CreateTier(newTier)
	if err != nil {
		return fmt.Errorf("ティアの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ティアを作成しました: %s (ID: %s)", createdTier.Name, createdTier.ID))
	}

	// JSON形式の場合はティア情報も出力
	if root.JSON {
		return formatter.Print(createdTier)
	}

	return nil
}

// TiersUpdateCmd はティアを更新するコマンドです
type TiersUpdateCmd struct {
	ID             string   `arg:"" help:"Tier ID"`
	Name           string   `help:"Tier name" short:"n"`
	Description    string   `help:"Tier description" short:"d"`
	Visibility     string   `help:"Visibility (public, none)"`
	MonthlyPrice   *int     `help:"Monthly price (in smallest currency unit)"`
	YearlyPrice    *int     `help:"Yearly price (in smallest currency unit)"`
	WelcomePageURL string   `help:"Welcome page URL"`
	Benefits       []string `help:"Benefits list" short:"b"`
}

// Run はtiersコマンドのupdateサブコマンドを実行します
func (c *TiersUpdateCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のティアを取得
	existingTier, err := client.GetTier(c.ID)
	if err != nil {
		return fmt.Errorf("ティアの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("update tier '%s' (ID: %s)", existingTier.Name, c.ID)
	if err := confirmDestructive(action, root.Force, root.NoInput); err != nil {
		return err
	}

	// 更新内容を反映
	updateTier := &ghostapi.Tier{
		Name:           existingTier.Name,
		Slug:           existingTier.Slug,
		Description:    existingTier.Description,
		Type:           existingTier.Type,
		Visibility:     existingTier.Visibility,
		MonthlyPrice:   existingTier.MonthlyPrice,
		YearlyPrice:    existingTier.YearlyPrice,
		Currency:       existingTier.Currency,
		WelcomePageURL: existingTier.WelcomePageURL,
		Benefits:       existingTier.Benefits,
	}

	if c.Name != "" {
		updateTier.Name = c.Name
	}
	if c.Description != "" {
		updateTier.Description = c.Description
	}
	if c.Visibility != "" {
		updateTier.Visibility = c.Visibility
	}
	if c.MonthlyPrice != nil {
		updateTier.MonthlyPrice = *c.MonthlyPrice
	}
	if c.YearlyPrice != nil {
		updateTier.YearlyPrice = *c.YearlyPrice
	}
	if c.WelcomePageURL != "" {
		updateTier.WelcomePageURL = c.WelcomePageURL
	}
	if len(c.Benefits) > 0 {
		updateTier.Benefits = c.Benefits
	}

	// ティアを更新
	updatedTier, err := client.UpdateTier(c.ID, updateTier)
	if err != nil {
		return fmt.Errorf("ティアの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("ティアを更新しました: %s (ID: %s)", updatedTier.Name, updatedTier.ID))
	}

	// JSON形式の場合はティア情報も出力
	if root.JSON {
		return formatter.Print(updatedTier)
	}

	return nil
}
