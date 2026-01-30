/**
 * tiers.go
 * ティア管理コマンド
 *
 * Ghostティアの閲覧機能を提供します。
 * ビジネス設定の誤変更リスクを回避するため、読み取り操作（List, Get）のみ実装しています。
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
	List TiersListCmd `cmd:"" help:"List tiers"`
	Get  TiersGetCmd  `cmd:"" help:"Get a tier"`
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
