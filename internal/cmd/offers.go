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

	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// OffersCmd はオファー管理コマンドです
type OffersCmd struct {
	List   OffersListCmd   `cmd:"" help:"List offers"`
	Get    OffersInfoCmd   `cmd:"" help:"Show offer information"`
	Create OffersCreateCmd `cmd:"" help:"Create an offer"`
	Update OffersUpdateCmd `cmd:"" help:"Update an offer"`

	// Phase 2: 状態変更
	Archive OffersArchiveCmd `cmd:"" help:"Archive an offer"`
}

// OffersListCmd is the command to retrieve オファー list
type OffersListCmd struct {
	Limit  int    `help:"Number of offers to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)" aliases:"where,w"`
}

// Run executes the list subcommand of the offers command
func (c *OffersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get offer list
	response, err := client.ListOffers(ghostapi.OfferListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
	})
	if err != nil {
		return fmt.Errorf("failed to list offers: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Offers)
	}

	// Output in table format
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

// OffersInfoCmd is the command to show オファー information
type OffersInfoCmd struct {
	ID string `arg:"" help:"Offer ID"`
}

// Run executes the info subcommand of the offers command
func (c *OffersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get offer
	offer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get offer: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(offer)
	}

	// Output in table format
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

// OffersCreateCmd is the command to create オファー
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

// Run executes the create subcommand of the offers command
func (c *OffersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Confirm destructive operation
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

	// Create new offer
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
		return fmt.Errorf("failed to create offer: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created offer: %s (ID: %s)", createdOffer.Name, createdOffer.ID))
	}

	// Also output offer information if JSON format
	if root.JSON {
		return formatter.Print(createdOffer)
	}

	return nil
}

// OffersUpdateCmd is the command to update オファー
type OffersUpdateCmd struct {
	ID                 string `arg:"" help:"Offer ID"`
	Name               string `help:"Offer name" short:"n"`
	DisplayTitle       string `help:"Display title" short:"t"`
	DisplayDescription string `help:"Display description" short:"d"`
	Amount             *int   `help:"Discount amount"`
	DurationInMonths   *int   `help:"Duration in months (for repeating)"`
}

// Run executes the update subcommand of the offers command
func (c *OffersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing offer
	existingOffer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get offer: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("update offer '%s' (ID: %s)", existingOffer.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Apply updates
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

	// Update offer
	updatedOffer, err := client.UpdateOffer(c.ID, updateOffer)
	if err != nil {
		return fmt.Errorf("failed to update offer: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated offer: %s (ID: %s)", updatedOffer.Name, updatedOffer.ID))
	}

	// Also output offer information if JSON format
	if root.JSON {
		return formatter.Print(updatedOffer)
	}

	return nil
}

// ========================================
// Phase 2: 状態変更
// ========================================

// OffersArchiveCmd is the command to archive オファー
type OffersArchiveCmd struct {
	ID string `arg:"" help:"Offer ID"`
}

// Run executes the archive subcommand of the offers command
func (c *OffersArchiveCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing offer
	existingOffer, err := client.GetOffer(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get offer: %w", err)
	}

	// Error if already archived
	if existingOffer.Status == "archived" {
		return fmt.Errorf("offer is already archived")
	}

	// Change status to archived
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

	// Update offer
	archivedOffer, err := client.UpdateOffer(c.ID, updateOffer)
	if err != nil {
		return fmt.Errorf("failed to archive offer: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("archived offer: %s (ID: %s)", archivedOffer.Name, archivedOffer.ID))
	}

	// Also output offer information if JSON format
	if root.JSON {
		return formatter.Print(archivedOffer)
	}

	return nil
}
