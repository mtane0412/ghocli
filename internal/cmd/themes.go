/**
 * themes.go
 * Theme management commands
 *
 * Provides functionality for listing, uploading, and activating Ghost themes.
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mtane0412/ghocli/internal/outfmt"
)

// ThemesCmd is the theme management command
type ThemesCmd struct {
	List     ThemesListCmd     `cmd:"" help:"List themes"`
	Upload   ThemesUploadCmd   `cmd:"" help:"Upload a theme"`
	Activate ThemesActivateCmd `cmd:"" help:"Activate a theme"`
	Delete   ThemesDeleteCmd   `cmd:"" help:"Delete a theme"`

	// Phase 3: Composite operations
	Install ThemesInstallCmd `cmd:"" help:"Upload and activate a theme"`
}

// ThemesListCmd is the command to retrieve theme list
type ThemesListCmd struct{}

// Run executes the list subcommand of the themes command
func (c *ThemesListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get theme list
	response, err := client.ListThemes()
	if err != nil {
		return fmt.Errorf("failed to list themes: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Themes)
	}

	// Output in table format
	headers := []string{"Name", "Active", "Version", "Description"}
	rows := make([][]string, len(response.Themes))
	for i, theme := range response.Themes {
		active := ""
		if theme.Active {
			active = "✓"
		}

		version := ""
		description := ""
		if theme.Package != nil {
			version = theme.Package.Version
			description = theme.Package.Description
		}

		rows[i] = []string{
			theme.Name,
			active,
			version,
			description,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// ThemesUploadCmd is the command to upload theme
type ThemesUploadCmd struct {
	File string `arg:"" help:"Theme zip file path" type:"existingfile"`
}

// Run executes the upload subcommand of the themes command
func (c *ThemesUploadCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Open file
	file, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get filename
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file information: %w", err)
	}

	// Upload theme
	theme, err := client.UploadTheme(file, fileInfo.Name())
	if err != nil {
		return fmt.Errorf("failed to upload theme: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("uploaded theme: %s", theme.Name))
		if theme.Package != nil && theme.Package.Version != "" {
			formatter.PrintMessage(fmt.Sprintf("version: %s", theme.Package.Version))
		}
	}

	// Also output theme information if JSON format
	if root.JSON {
		return formatter.Print(theme)
	}

	return nil
}

// ThemesActivateCmd is the command to activate theme
type ThemesActivateCmd struct {
	Name string `arg:"" help:"Theme name"`
}

// Run executes the activate subcommand of the themes command
func (c *ThemesActivateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Activate theme
	theme, err := client.ActivateTheme(c.Name)
	if err != nil {
		return fmt.Errorf("failed to activate theme: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("activated theme: %s", theme.Name))
	}

	// Also output theme information if JSON format
	if root.JSON {
		return formatter.Print(theme)
	}

	return nil
}

// ========================================
// Phase 3: Composite operations
// ========================================

// ThemesInstallCmd is the command to upload and activate theme
type ThemesInstallCmd struct {
	File string `arg:"" help:"Path to theme zip file" type:"existingfile"`
}

// Run executes the install subcommand of the themes command
func (c *ThemesInstallCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Open file
	file, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get filename
	filename := filepath.Base(c.File)

	// Upload theme
	formatter.PrintMessage(fmt.Sprintf("uploading theme: %s", c.File))
	uploadedTheme, err := client.UploadTheme(file, filename)
	if err != nil {
		return fmt.Errorf("failed to upload theme: %w", err)
	}

	formatter.PrintMessage(fmt.Sprintf("uploaded theme: %s", uploadedTheme.Name))

	// Activate theme
	formatter.PrintMessage(fmt.Sprintf("activating theme: %s", uploadedTheme.Name))
	activatedTheme, err := client.ActivateTheme(uploadedTheme.Name)
	if err != nil {
		return fmt.Errorf("failed to activate theme: %w", err)
	}

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("installed and activated theme: %s", activatedTheme.Name))
	}

	// Also output theme information if JSON format
	if root.JSON {
		return formatter.Print(activatedTheme)
	}

	return nil
}

// ========================================
// Theme deletion
// ========================================

// ThemesDeleteCmd is the command to delete theme
type ThemesDeleteCmd struct {
	Name string `arg:"" help:"Theme name"`
}

// Run executes the delete subcommand of the themes command
func (c *ThemesDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get theme list and check if active
	themes, err := client.ListThemes()
	if err != nil {
		return fmt.Errorf("failed to list themes: %w", err)
	}

	// Prevent deletion of active theme
	for _, theme := range themes.Themes {
		if theme.Name == c.Name && theme.Active {
			return fmt.Errorf("cannot delete active theme: %s", c.Name)
		}
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete theme '%s'", c.Name)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Delete theme
	if err := client.DeleteTheme(c.Name); err != nil {
		return fmt.Errorf("failed to delete theme: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted theme: %s", c.Name))

	return nil
}
