/**
 * settings.go
 * 設定管理コマンド
 *
 * Ghostサイト設定の表示、更新機能を提供します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// SettingsCmd は設定管理コマンドです
type SettingsCmd struct {
	List SettingsListCmd `cmd:"" help:"List all settings"`
	Get  SettingsGetCmd  `cmd:"" help:"Get a specific setting"`
	Set  SettingsSetCmd  `cmd:"" help:"Set a setting value"`
}

// SettingsListCmd is the command to retrieve 設定 list
type SettingsListCmd struct{}

// Run executes the list subcommand of the settings command
func (c *SettingsListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get settings list
	response, err := client.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Settings)
	}

	// Output in table format
	headers := []string{"Key", "Value"}
	rows := make([][]string, len(response.Settings))
	for i, setting := range response.Settings {
		value := fmt.Sprintf("%v", setting.Value)
		// Truncate long values
		if len(value) > 80 {
			value = value[:77] + "..."
		}
		rows[i] = []string{
			setting.Key,
			value,
		}
	}

	return formatter.PrintTable(headers, rows)
}

// SettingsGetCmd is the command to get specific setting value
type SettingsGetCmd struct {
	Key string `arg:"" help:"Setting key"`
}

// Run executes the get subcommand of the settings command
func (c *SettingsGetCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get settings list
	response, err := client.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Search for settings with specified key
	var foundSetting *ghostapi.Setting
	for _, setting := range response.Settings {
		if setting.Key == c.Key {
			foundSetting = &setting
			break
		}
	}

	if foundSetting == nil {
		return fmt.Errorf("setting not found: %s", c.Key)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(foundSetting)
	}

	// Output value
	formatter.PrintMessage(fmt.Sprintf("%v", foundSetting.Value))

	return nil
}

// SettingsSetCmd is the command to update 設定 value
type SettingsSetCmd struct {
	Key   string `arg:"" help:"Setting key"`
	Value string `arg:"" help:"Setting value"`
}

// Run executes the set subcommand of the settings command
func (c *SettingsSetCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Update settings
	updates := []ghostapi.SettingUpdate{
		{
			Key:   c.Key,
			Value: c.Value,
		},
	}

	response, err := client.UpdateSettings(updates)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated settings: %s = %s", c.Key, c.Value))
	}

	// Also output settings information if JSON format
	if root.JSON {
		return formatter.Print(response.Settings)
	}

	return nil
}
