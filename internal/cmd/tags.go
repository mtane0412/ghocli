/**
 * tags.go
 * タグ管理コマンド
 *
 * Ghostタグの作成、更新、削除機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/fields"
	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// TagsCmd はタグ管理コマンドです
type TagsCmd struct {
	List   TagsListCmd   `cmd:"" help:"List tags"`
	Get    TagsInfoCmd   `cmd:"" help:"Show tag information"`
	Create TagsCreateCmd `cmd:"" help:"Create a tag"`
	Update TagsUpdateCmd `cmd:"" help:"Update a tag"`
	Delete TagsDeleteCmd `cmd:"" help:"Delete a tag"`
}

// TagsListCmd is the command to retrieve タグ list
type TagsListCmd struct {
	Limit   int    `help:"Number of tags to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (count.posts)" short:"i"`
}

// Run executes the list subcommand of the tags command
func (c *TagsListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.TagFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get tag list
	response, err := client.ListTags(ghostapi.TagListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
	})
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Tag struct to map[string]interface{}
		var tagsData []map[string]interface{}
		for _, tag := range response.Tags {
			tagMap, err := outfmt.StructToMap(tag)
			if err != nil {
				return fmt.Errorf("failed to convert tag data: %w", err)
			}
			tagsData = append(tagsData, tagMap)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, tagsData, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Tags)
	}

	// Output in table format
	headers := []string{"ID", "Name", "Slug", "Visibility", "Created"}
	rows := make([][]string, len(response.Tags))
	for i, tag := range response.Tags {
		rows[i] = []string{
			tag.ID,
			tag.Name,
			tag.Slug,
			tag.Visibility,
			tag.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// TagsInfoCmd is the command to show タグ information
type TagsInfoCmd struct {
	IDOrSlug string `arg:"" help:"Tag ID or slug (use 'slug:tag-name' format for slug)"`
}

// Run executes the info subcommand of the tags command
func (c *TagsInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.TagFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get tag
	tag, err := client.GetTag(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Tag struct to map[string]interface{}
		tagMap, err := outfmt.StructToMap(tag)
		if err != nil {
			return fmt.Errorf("failed to convert tag data: %w", err)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, []map[string]interface{}{tagMap}, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(tag)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"id", tag.ID},
		{"name", tag.Name},
		{"slug", tag.Slug},
		{"description", tag.Description},
		{"visibility", tag.Visibility},
		{"created", tag.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", tag.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// TagsCreateCmd is the command to create タグ
type TagsCreateCmd struct {
	Name        string `help:"Tag name" short:"n" required:""`
	Description string `help:"Tag description" short:"d"`
	Visibility  string `help:"Tag visibility (public, internal)" default:"public"`
}

// Run executes the create subcommand of the tags command
func (c *TagsCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create new tag
	newTag := &ghostapi.Tag{
		Name:        c.Name,
		Description: c.Description,
		Visibility:  c.Visibility,
	}

	createdTag, err := client.CreateTag(newTag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created tag: %s (ID: %s)", createdTag.Name, createdTag.ID))
	}

	// Also output tag information if JSON format
	if root.JSON {
		return formatter.Print(createdTag)
	}

	return nil
}

// TagsUpdateCmd is the command to update タグ
type TagsUpdateCmd struct {
	ID          string `arg:"" help:"Tag ID"`
	Name        string `help:"Tag name" short:"n"`
	Description string `help:"Tag description" short:"d"`
	Visibility  string `help:"Tag visibility (public, internal)"`
}

// Run executes the update subcommand of the tags command
func (c *TagsUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing tag
	existingTag, err := client.GetTag(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Apply updates
	updateTag := &ghostapi.Tag{
		Name:        existingTag.Name,
		Slug:        existingTag.Slug,
		Description: existingTag.Description,
		Visibility:  existingTag.Visibility,
	}

	if c.Name != "" {
		updateTag.Name = c.Name
	}
	if c.Description != "" {
		updateTag.Description = c.Description
	}
	if c.Visibility != "" {
		updateTag.Visibility = c.Visibility
	}

	// Update tag
	updatedTag, err := client.UpdateTag(c.ID, updateTag)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated tag: %s (ID: %s)", updatedTag.Name, updatedTag.ID))
	}

	// Also output tag information if JSON format
	if root.JSON {
		return formatter.Print(updatedTag)
	}

	return nil
}

// TagsDeleteCmd is the command to delete タグ
type TagsDeleteCmd struct {
	ID string `arg:"" help:"Tag ID"`
}

// Run executes the delete subcommand of the tags command
func (c *TagsDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get tag information to build confirmation message
	tag, err := client.GetTag(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete tag '%s' (ID: %s)", tag.Name, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Delete tag
	if err := client.DeleteTag(c.ID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted tag (ID: %s)", c.ID))

	return nil
}
