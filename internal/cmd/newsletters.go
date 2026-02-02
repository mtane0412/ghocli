/**
 * newsletters.go
 * ニュースレター管理コマンド
 *
 * Ghostニュースレターの管理機能を提供します。
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

// NewslettersCmd はニュースレター管理コマンドです
type NewslettersCmd struct {
	List   NewslettersListCmd   `cmd:"" help:"List newsletters"`
	Get    NewslettersInfoCmd   `cmd:"" help:"Show newsletter information"`
	Create NewslettersCreateCmd `cmd:"" help:"Create a newsletter"`
	Update NewslettersUpdateCmd `cmd:"" help:"Update a newsletter"`
}

// NewslettersListCmd is the command to retrieve ニュースレター list
type NewslettersListCmd struct {
	Limit  int    `help:"Number of newsletters to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter condition (e.g., status:active)" aliases:"where,w"`
}

// Run executes the list subcommand of the newsletters command
func (c *NewslettersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get newsletter list
	response, err := client.ListNewsletters(ghostapi.NewsletterListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
	})
	if err != nil {
		return fmt.Errorf("failed to list newsletters: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Newsletters)
	}

	// Output in table format
	headers := []string{"ID", "Name", "Slug", "Status", "Visibility", "Created"}
	rows := make([][]string, len(response.Newsletters))
	for i, newsletter := range response.Newsletters {
		rows[i] = []string{
			newsletter.ID,
			newsletter.Name,
			newsletter.Slug,
			newsletter.Status,
			newsletter.Visibility,
			newsletter.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// NewslettersInfoCmd is the command to show ニュースレター information
type NewslettersInfoCmd struct {
	IDOrSlug string `arg:"" help:"Newsletter ID or slug (use 'slug:newsletter-name' format for slug)"`
}

// Run executes the info subcommand of the newsletters command
func (c *NewslettersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get newsletter
	newsletter, err := client.GetNewsletter(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get newsletter: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(newsletter)
	}

	// Output in table format
	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"ID", newsletter.ID},
		{"Name", newsletter.Name},
		{"Slug", newsletter.Slug},
		{"Description", newsletter.Description},
		{"Status", newsletter.Status},
		{"Visibility", newsletter.Visibility},
		{"Subscribe on Signup", fmt.Sprintf("%t", newsletter.SubscribeOnSignup)},
		{"Sender Name", newsletter.SenderName},
		{"Sender Email", newsletter.SenderEmail},
		{"Sender Reply To", newsletter.SenderReplyTo},
		{"Sort Order", fmt.Sprintf("%d", newsletter.SortOrder)},
		{"Created", newsletter.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", newsletter.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	return formatter.PrintTable(headers, rows)
}

// NewslettersCreateCmd is the command to create ニュースレター
type NewslettersCreateCmd struct {
	Name              string `help:"Newsletter name" short:"n" required:""`
	Description       string `help:"Newsletter description" short:"d"`
	Visibility        string `help:"Visibility (members, paid)" default:"members"`
	SubscribeOnSignup bool   `help:"Subscribe members on signup" default:"true"`
	SenderName        string `help:"Sender name"`
	SenderEmail       string `help:"Sender email"`
}

// Run executes the create subcommand of the newsletters command
func (c *NewslettersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Confirm destructive operation
	action := fmt.Sprintf("create newsletter '%s'", c.Name)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Create new newsletter
	newNewsletter := &ghostapi.Newsletter{
		Name:              c.Name,
		Description:       c.Description,
		Visibility:        c.Visibility,
		SubscribeOnSignup: c.SubscribeOnSignup,
		SenderName:        c.SenderName,
		SenderEmail:       c.SenderEmail,
	}

	createdNewsletter, err := client.CreateNewsletter(newNewsletter)
	if err != nil {
		return fmt.Errorf("failed to create newsletter: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created newsletter: %s (ID: %s)", createdNewsletter.Name, createdNewsletter.ID))
	}

	// Also output newsletter information if JSON format
	if root.JSON {
		return formatter.Print(createdNewsletter)
	}

	return nil
}

// NewslettersUpdateCmd is the command to update ニュースレター
type NewslettersUpdateCmd struct {
	ID                string `arg:"" help:"Newsletter ID"`
	Name              string `help:"Newsletter name" short:"n"`
	Description       string `help:"Newsletter description" short:"d"`
	Visibility        string `help:"Visibility (members, paid)"`
	SubscribeOnSignup *bool  `help:"Subscribe members on signup"`
	SenderName        string `help:"Sender name"`
	SenderEmail       string `help:"Sender email"`
}

// Run executes the update subcommand of the newsletters command
func (c *NewslettersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing newsletter
	existingNewsletter, err := client.GetNewsletter(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get newsletter: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("update newsletter '%s' (ID: %s)", existingNewsletter.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Apply updates
	updateNewsletter := &ghostapi.Newsletter{
		Name:              existingNewsletter.Name,
		Slug:              existingNewsletter.Slug,
		Description:       existingNewsletter.Description,
		Visibility:        existingNewsletter.Visibility,
		SubscribeOnSignup: existingNewsletter.SubscribeOnSignup,
		SenderName:        existingNewsletter.SenderName,
		SenderEmail:       existingNewsletter.SenderEmail,
		SenderReplyTo:     existingNewsletter.SenderReplyTo,
	}

	if c.Name != "" {
		updateNewsletter.Name = c.Name
	}
	if c.Description != "" {
		updateNewsletter.Description = c.Description
	}
	if c.Visibility != "" {
		updateNewsletter.Visibility = c.Visibility
	}
	if c.SubscribeOnSignup != nil {
		updateNewsletter.SubscribeOnSignup = *c.SubscribeOnSignup
	}
	if c.SenderName != "" {
		updateNewsletter.SenderName = c.SenderName
	}
	if c.SenderEmail != "" {
		updateNewsletter.SenderEmail = c.SenderEmail
	}

	// Update newsletter
	updatedNewsletter, err := client.UpdateNewsletter(c.ID, updateNewsletter)
	if err != nil {
		return fmt.Errorf("failed to update newsletter: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated newsletter: %s (ID: %s)", updatedNewsletter.Name, updatedNewsletter.ID))
	}

	// Also output newsletter information if JSON format
	if root.JSON {
		return formatter.Print(updatedNewsletter)
	}

	return nil
}
