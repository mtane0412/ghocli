/**
 * tiers.go
 * ティア管理コマンド
 *
 * Ghostティアの管理機能を提供します。
 * Create/Update操作には確認機構が適用されます。
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// TiersCmd はティア管理コマンドです
type TiersCmd struct {
	List   TiersListCmd   `cmd:"" help:"List tiers"`
	Get    TiersInfoCmd   `cmd:"" help:"Show tier information"`
	Create TiersCreateCmd `cmd:"" help:"Create a tier"`
	Update TiersUpdateCmd `cmd:"" help:"Update a tier"`
}

// TiersListCmd is the command to retrieve ティア list
type TiersListCmd struct {
	Limit   int    `help:"Number of tiers to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (monthly_price,yearly_price,benefits)" short:"i"`
	Filter  string `help:"Filter condition" aliases:"where,w"`
}

// Run executes the list subcommand of the tiers command
func (c *TiersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get tier list
	response, err := client.ListTiers(ghostapi.TierListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
		Filter:  c.Filter,
	})
	if err != nil {
		return fmt.Errorf("failed to list tiers: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Tiers)
	}

	// Output in table format
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

// TiersInfoCmd is the command to show ティア information
type TiersInfoCmd struct {
	IDOrSlug string `arg:"" help:"Tier ID or slug (use 'slug:tier-name' format for slug)"`
	Include  string `help:"Include additional data (monthly_price,yearly_price,benefits)" short:"i"`
}

// Run executes the info subcommand of the tiers command
func (c *TiersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get tier
	tier, err := client.GetTier(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get tier: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(tier)
	}

	// Output in table format
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

// TiersCreateCmd is the command to create ティア
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

// Run executes the create subcommand of the tiers command
func (c *TiersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Confirm destructive operation
	priceInfo := ""
	if c.Type == "paid" {
		priceInfo = fmt.Sprintf(" (monthly: %d %s, yearly: %d %s)", c.MonthlyPrice, c.Currency, c.YearlyPrice, c.Currency)
	}
	action := fmt.Sprintf("create tier '%s'%s", c.Name, priceInfo)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Create new tier
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
		return fmt.Errorf("failed to create tier: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created tier: %s (ID: %s)", createdTier.Name, createdTier.ID))
	}

	// Also output tier information if JSON format
	if root.JSON {
		return formatter.Print(createdTier)
	}

	return nil
}

// TiersUpdateCmd is the command to update ティア
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

// Run executes the update subcommand of the tiers command
func (c *TiersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing tier
	existingTier, err := client.GetTier(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get tier: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("update tier '%s' (ID: %s)", existingTier.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Apply updates
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

	// Update tier
	updatedTier, err := client.UpdateTier(c.ID, updateTier)
	if err != nil {
		return fmt.Errorf("failed to update tier: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated tier: %s (ID: %s)", updatedTier.Name, updatedTier.ID))
	}

	// Also output tier information if JSON format
	if root.JSON {
		return formatter.Print(updatedTier)
	}

	return nil
}
