/**
 * config.go
 * Configuration management commands
 *
 * Provides commands to manage gho configuration.
 */

package cmd

import (
	"context"
	"fmt"

	"github.com/mtane0412/gho/internal/config"
)

// ConfigCmd is the root command for configuration management
type ConfigCmd struct {
	Get    ConfigGetCmd    `cmd:"" help:"Get configuration value"`
	Set    ConfigSetCmd    `cmd:"" help:"Set configuration value"`
	Unset  ConfigUnsetCmd  `cmd:"" help:"Unset configuration value"`
	List   ConfigListCmd   `cmd:"" help:"List all configuration"`
	Path   ConfigPathCmd   `cmd:"" help:"Show configuration file path"`
	Keys   ConfigKeysCmd   `cmd:"" help:"List available configuration keys"`
}

// ConfigGetCmd is the command to get configuration values
type ConfigGetCmd struct {
	Key string `arg:"" help:"Configuration key to get"`
}

// Run executes the config get command
func (c *ConfigGetCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get value based on key
	var value string
	switch c.Key {
	case "default_site":
		value = cfg.DefaultSite
	case "keyring_backend":
		value = cfg.KeyringBackend
	default:
		return fmt.Errorf("unknown configuration key: %s", c.Key)
	}

	// Output value
	fmt.Println(value)
	return nil
}

// ConfigSetCmd is the command to set configuration values
type ConfigSetCmd struct {
	Key   string `arg:"" help:"Configuration key to set"`
	Value string `arg:"" help:"Configuration value to set"`
}

// Run executes the config set command
func (c *ConfigSetCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Set value based on key
	switch c.Key {
	case "default_site":
		cfg.DefaultSite = c.Value
	case "keyring_backend":
		cfg.KeyringBackend = c.Value
	default:
		return fmt.Errorf("unknown configuration key: %s", c.Key)
	}

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Set %s to %s\n", c.Key, c.Value)
	return nil
}

// ConfigUnsetCmd is the command to unset configuration values
type ConfigUnsetCmd struct {
	Key string `arg:"" help:"Configuration key to unset"`
}

// Run executes the config unset command
func (c *ConfigUnsetCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Unset value (set to empty string) based on key
	switch c.Key {
	case "default_site":
		cfg.DefaultSite = ""
	case "keyring_backend":
		cfg.KeyringBackend = ""
	default:
		return fmt.Errorf("unknown configuration key: %s", c.Key)
	}

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Unset %s\n", c.Key)
	return nil
}

// ConfigListCmd is the command to list all configuration
type ConfigListCmd struct{}

// Run executes the config list command
func (c *ConfigListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Display configuration
	fmt.Printf("default_site=%s\n", cfg.DefaultSite)
	fmt.Printf("keyring_backend=%s\n", cfg.KeyringBackend)

	return nil
}

// ConfigPathCmd is the command to display configuration file path
type ConfigPathCmd struct{}

// Run executes the config path command
func (c *ConfigPathCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	fmt.Println(configPath)
	return nil
}

// ConfigKeysCmd is the command to list available configuration keys
type ConfigKeysCmd struct{}

// Run executes the config keys command
func (c *ConfigKeysCmd) Run(ctx context.Context, root *RootFlags) error {
	fmt.Println("default_site")
	fmt.Println("keyring_backend")
	return nil
}
