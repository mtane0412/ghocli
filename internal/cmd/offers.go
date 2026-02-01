/**
 * offers.go
 * オファー管理コマンド
 *
 * Ghostオファーの管理機能を提供します。
 * Create/Update操作には確認機構が適用されます。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
)

// OffersCmd はオファー管理コマンドです
type OffersCmd struct {
	List   OffersListCmd   `cmd:"" help:"List offers"`
	Get    OffersInfoCmd   `cmd:"" help:"オファーの情報を表示"`
	Create OffersCreateCmd `cmd:"" help:"Create an offer"`
	Update OffersUpdateCmd `cmd:"" help:"Update an offer"`

	// Phase 2: 状態変更
	Archive OffersArchiveCmd `cmd:"" help:"Archive an offer"`
}

// OffersListCmd はオファー一覧を取得するコマンドです
type OffersListCmd struct {
	Limit  int    `help:"Number of offers to retrieve" short:"l" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)"`
}

// Run はoffersコマンドのlistサブコマンドを実行します
func (c *OffersListCmd) Run(ctx context.Context, root *RootFlags) error {
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

// OffersInfoCmd はオファー情報を表示するコマンドです
type OffersInfoCmd struct {
	ID string `arg:"" help:"Offer ID"`
}

// Run はoffersコマンドのinfoサブコマンドを実行します
func (c *OffersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
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

// OffersCreateCmd はオファーを作成するコマンドです
type OffersCreateCmd struct {
	Name               string `help:"Offer name" short:"n" required:""`
	Code               string `help:"Offer code" short:"c" required:""`
	DisplayTitle       string `help:"Display title" short:"t"`
	DisplayDescription string `help:"Display description" short:"d"`
	Type               string `help:"Offer type (percent, fixed)" default:"percent"`
	Cadence            string `help:"Cadence (month, year)" default:"month"`
	Amount             int    `help:"Discount amount" required:""`
	Duration           string `help:"Duration (once, forever, repeating)" default:"once"`
	DurationInMonths   int    `help:"Duration in months (for repeating)"`
	Currency           string `help:"Currency code (for fixed type)" default:"JPY"`
	TierID             string `help:"Tier ID" required:""`
}

// Run はoffersコマンドのcreateサブコマンドを実行します
func (c *OffersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 破壊的操作の確認
	discountInfo := fmt.Sprintf("%d", c.Amount)
	if c.Type == "percent" {
		discountInfo += "%"
	} else {
		discountInfo += " " + c.Currency
	}
	action := fmt.Sprintf("create offer '%s' (code: %s, discount: %s)", c.Name, c.Code, discountInfo)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// 新規オファーを作成
	newOffer := &ghostapi.Offer{
		Name:               c.Name,
		Code:               c.Code,
		DisplayTitle:       c.DisplayTitle,
		DisplayDescription: c.DisplayDescription,
		Type:               c.Type,
		Cadence:            c.Cadence,
		Amount:             c.Amount,
		Duration:           c.Duration,
		DurationInMonths:   c.DurationInMonths,
		Currency:           c.Currency,
		Tier: ghostapi.OfferTier{
			ID: c.TierID,
		},
	}

	createdOffer, err := client.CreateOffer(newOffer)
	if err != nil {
		return fmt.Errorf("オファーの作成に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("オファーを作成しました: %s (ID: %s)", createdOffer.Name, createdOffer.ID))
	}

	// JSON形式の場合はオファー情報も出力
	if root.JSON {
		return formatter.Print(createdOffer)
	}

	return nil
}

// OffersUpdateCmd はオファーを更新するコマンドです
type OffersUpdateCmd struct {
	ID                 string `arg:"" help:"Offer ID"`
	Name               string `help:"Offer name" short:"n"`
	DisplayTitle       string `help:"Display title" short:"t"`
	DisplayDescription string `help:"Display description" short:"d"`
	Amount             *int   `help:"Discount amount"`
	DurationInMonths   *int   `help:"Duration in months (for repeating)"`
}

// Run はoffersコマンドのupdateサブコマンドを実行します
func (c *OffersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のオファーを取得
	existingOffer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("オファーの取得に失敗: %w", err)
	}

	// 破壊的操作の確認
	action := fmt.Sprintf("update offer '%s' (ID: %s)", existingOffer.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// 更新内容を反映
	updateOffer := &ghostapi.Offer{
		Name:               existingOffer.Name,
		Code:               existingOffer.Code,
		DisplayTitle:       existingOffer.DisplayTitle,
		DisplayDescription: existingOffer.DisplayDescription,
		Type:               existingOffer.Type,
		Cadence:            existingOffer.Cadence,
		Amount:             existingOffer.Amount,
		Duration:           existingOffer.Duration,
		DurationInMonths:   existingOffer.DurationInMonths,
		Currency:           existingOffer.Currency,
		Tier:               existingOffer.Tier,
	}

	if c.Name != "" {
		updateOffer.Name = c.Name
	}
	if c.DisplayTitle != "" {
		updateOffer.DisplayTitle = c.DisplayTitle
	}
	if c.DisplayDescription != "" {
		updateOffer.DisplayDescription = c.DisplayDescription
	}
	if c.Amount != nil {
		updateOffer.Amount = *c.Amount
	}
	if c.DurationInMonths != nil {
		updateOffer.DurationInMonths = *c.DurationInMonths
	}

	// オファーを更新
	updatedOffer, err := client.UpdateOffer(c.ID, updateOffer)
	if err != nil {
		return fmt.Errorf("オファーの更新に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("オファーを更新しました: %s (ID: %s)", updatedOffer.Name, updatedOffer.ID))
	}

	// JSON形式の場合はオファー情報も出力
	if root.JSON {
		return formatter.Print(updatedOffer)
	}

	return nil
}

// ========================================
// Phase 2: 状態変更
// ========================================

// OffersArchiveCmd はオファーをアーカイブするコマンドです
type OffersArchiveCmd struct {
	ID string `arg:"" help:"Offer ID"`
}

// Run はoffersコマンドのarchiveサブコマンドを実行します
func (c *OffersArchiveCmd) Run(ctx context.Context, root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// 既存のオファーを取得
	existingOffer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("オファーの取得に失敗: %w", err)
	}

	// すでにアーカイブ済みの場合はエラー
	if existingOffer.Status == "archived" {
		return fmt.Errorf("このオファーはすでにアーカイブされています")
	}

	// ステータスをarchivedに変更
	updateOffer := &ghostapi.Offer{
		Name:               existingOffer.Name,
		Code:               existingOffer.Code,
		DisplayTitle:       existingOffer.DisplayTitle,
		DisplayDescription: existingOffer.DisplayDescription,
		Type:               existingOffer.Type,
		Cadence:            existingOffer.Cadence,
		Amount:             existingOffer.Amount,
		Duration:           existingOffer.Duration,
		DurationInMonths:   existingOffer.DurationInMonths,
		Currency:           existingOffer.Currency,
		Status:             "archived",
		Tier:               existingOffer.Tier,
	}

	// オファーを更新
	archivedOffer, err := client.UpdateOffer(c.ID, updateOffer)
	if err != nil {
		return fmt.Errorf("オファーのアーカイブに失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// 成功メッセージを表示
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("オファーをアーカイブしました: %s (ID: %s)", archivedOffer.Name, archivedOffer.ID))
	}

	// JSON形式の場合はオファー情報も出力
	if root.JSON {
		return formatter.Print(archivedOffer)
	}

	return nil
}
