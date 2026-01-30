/**
 * offers.go
 * オファー管理コマンド
 *
 * Ghostオファーの閲覧機能を提供します。
 * ビジネス設定の誤変更リスクを回避するため、読み取り操作（List, Get）のみ実装しています。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// OffersCmd はオファー管理コマンドです
type OffersCmd struct {
	List OffersListCmd `cmd:"" help:"List offers"`
	Get  OffersGetCmd  `cmd:"" help:"Get an offer"`
}

// OffersListCmd はオファー一覧を取得するコマンドです
type OffersListCmd struct {
	Limit  int    `help:"Number of offers to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)"`
}

// Run はoffersコマンドのlistサブコマンドを実行します
func (c *OffersListCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// オファー一覧を取得
	response, err := client.ListOffers(ghostapi.OfferListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
	})
	if err != nil {
		return fmt.Errorf("オファー一覧の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(response.Offers)
	}

	// テーブル形式で出力
	headers := []string{"ID", "Name", "Code", "Type", "Amount", "Status", "Redemptions", "Created"}
	rows := make([][]string, len(response.Offers))
	for i, offer := range response.Offers {
		rows[i] = []string{
			offer.ID,
			offer.Name,
			offer.Code,
			offer.Type,
			fmt.Sprintf("%d", offer.Amount),
			offer.Status,
			fmt.Sprintf("%d", offer.RedemptionCount),
			offer.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// OffersGetCmd はオファーを取得するコマンドです
type OffersGetCmd struct {
	ID string `arg:"" help:"Offer ID"`
}

// Run はoffersコマンドのgetサブコマンドを実行します
func (c *OffersGetCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// オファーを取得
	offer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("オファーの取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(offer)
	}

	// テーブル形式で出力
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", offer.ID},
		{"Name", offer.Name},
		{"Code", offer.Code},
		{"Display Title", offer.DisplayTitle},
		{"Display Description", offer.DisplayDescription},
		{"Type", offer.Type},
		{"Cadence", offer.Cadence},
		{"Amount", fmt.Sprintf("%d", offer.Amount)},
		{"Duration", offer.Duration},
		{"Duration in Months", fmt.Sprintf("%d", offer.DurationInMonths)},
		{"Currency", offer.Currency},
		{"Status", offer.Status},
		{"Redemption Count", fmt.Sprintf("%d", offer.RedemptionCount)},
		{"Tier ID", offer.Tier.ID},
		{"Tier Name", offer.Tier.Name},
		{"Created", offer.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", offer.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	return formatter.PrintTable(headers, rows)
}
