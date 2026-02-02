/**
 * webhooks.go
 * Webhook management commands
 *
 * Provides functionality for creating, updating, and deleting Ghost webhooks.
 * Note: Ghost API does not support List/Get operations for webhooks.
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// WebhooksCmd is the webhook management command
type WebhooksCmd struct {
	Create WebhooksCreateCmd `cmd:"" help:"Create a webhook"`
	Update WebhooksUpdateCmd `cmd:"" help:"Update a webhook"`
	Delete WebhooksDeleteCmd `cmd:"" help:"Delete a webhook"`
}

// WebhooksCreateCmd is the command to create Webhook
type WebhooksCreateCmd struct {
	Event     string `help:"Webhook event (e.g., post.published, member.added)" short:"e" required:""`
	TargetURL string `help:"Target URL for webhook" short:"t" required:""`
	Name      string `help:"Webhook name" short:"n"`
}

// Run executes the create subcommand of the webhooks command
func (c *WebhooksCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create new webhook
	webhook := &ghostapi.Webhook{
		Event:     c.Event,
		TargetURL: c.TargetURL,
		Name:      c.Name,
	}

	created, err := client.CreateWebhook(webhook)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created webhook (ID: %s)", created.ID))
		formatter.PrintMessage(fmt.Sprintf("event: %s", created.Event))
		formatter.PrintMessage(fmt.Sprintf("URL: %s", created.TargetURL))
		if created.Secret != "" {
			formatter.PrintMessage(fmt.Sprintf("secret: %s", created.Secret))
		}
	}

	// Also output webhook information if JSON format
	if root.JSON {
		return formatter.Print(created)
	}

	return nil
}

// WebhooksUpdateCmd is the command to update Webhook
type WebhooksUpdateCmd struct {
	ID        string `arg:"" help:"Webhook ID"`
	Event     string `help:"Webhook event" short:"e"`
	TargetURL string `help:"Target URL for webhook" short:"t"`
	Name      string `help:"Webhook name" short:"n"`
}

// Run executes the update subcommand of the webhooks command
func (c *WebhooksUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create updates (specified fields only)
	webhook := &ghostapi.Webhook{}

	if c.Event != "" {
		webhook.Event = c.Event
	}
	if c.TargetURL != "" {
		webhook.TargetURL = c.TargetURL
	}
	if c.Name != "" {
		webhook.Name = c.Name
	}

	// Update webhook
	updated, err := client.UpdateWebhook(c.ID, webhook)
	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated webhook (ID: %s)", updated.ID))
		formatter.PrintMessage(fmt.Sprintf("event: %s", updated.Event))
		formatter.PrintMessage(fmt.Sprintf("URL: %s", updated.TargetURL))
	}

	// Also output webhook information if JSON format
	if root.JSON {
		return formatter.Print(updated)
	}

	return nil
}

// WebhooksDeleteCmd is the command to delete Webhook
type WebhooksDeleteCmd struct {
	ID string `arg:"" help:"Webhook ID"`
}

// Run executes the delete subcommand of the webhooks command
func (c *WebhooksDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Request confirmation unless skipping confirmation
	if !root.Force {
		fmt.Printf("Are you sure you want to delete Webhook (ID: %s)? [y/N]: ", c.ID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Delete webhook
	if err := client.DeleteWebhook(c.ID); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted webhook (ID: %s)", c.ID))

	return nil
}
