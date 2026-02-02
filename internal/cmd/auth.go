/**
 * auth.go
 * Authentication management commands
 *
 * Performs addition, listing, removal, and status checking of Ghost Admin API keys.
 */

package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mtane0412/gho/internal/config"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
	"github.com/mtane0412/gho/internal/secrets"
)

// AuthCmd is the root command for authentication management
type AuthCmd struct {
	Add         AuthAddCmd         `cmd:"" help:"Add a new site authentication"`
	List        AuthListCmd        `cmd:"" help:"List authenticated sites"`
	Remove      AuthRemoveCmd      `cmd:"" help:"Remove site authentication"`
	Status      AuthStatusCmd      `cmd:"" help:"Check authentication status"`
	Tokens      AuthTokensCmd      `cmd:"" help:"Manage API tokens"`
	Credentials AuthCredentialsCmd `cmd:"" help:"Add authentication from credentials file"`
}

// AuthAddCmd is the command to add site authentication
type AuthAddCmd struct {
	SiteURL string `arg:"" help:"Ghost site URL (e.g., https://myblog.ghost.io)"`
	Alias   string `help:"Alias for this site" short:"a"`
}

// Run executes the auth add command
func (c *AuthAddCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// Prompt for API key input
	fmt.Print("Enter Admin API Key (id:secret): ")
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)

	// Parse API key
	keyID, secret, err := secrets.ParseAdminAPIKey(apiKey)
	if err != nil {
		return err
	}

	// Validate API key (by fetching site info)
	client, err := ghostapi.NewClient(c.SiteURL, keyID, secret)
	if err != nil {
		return err
	}

	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}

	// Determine alias
	alias := c.Alias
	if alias == "" {
		// Generate alias from URL (e.g., https://myblog.ghost.io -> myblog)
		alias = extractAliasFromURL(c.SiteURL)
	}

	// Save to keyring
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	if err := store.Set(alias, apiKey); err != nil {
		return err
	}

	// Add site to configuration
	cfg.AddSite(alias, c.SiteURL)
	if cfg.DefaultSite == "" {
		cfg.DefaultSite = alias
	}

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Added site '%s' (%s)\n", alias, site.Title)
	return nil
}

// AuthListCmd is the command to list authenticated sites
type AuthListCmd struct{}

// Run executes the auth list command
func (c *AuthListCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// Get site list
	if len(cfg.Sites) == 0 {
		fmt.Println("No sites configured")
		return nil
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output in table format
	headers := []string{"Alias", "URL", "Default"}
	var rows [][]string
	for alias, url := range cfg.Sites {
		isDefault := ""
		if alias == cfg.DefaultSite {
			isDefault = "*"
		}
		rows = append(rows, []string{alias, url, isDefault})
	}

	return formatter.PrintTable(headers, rows)
}

// AuthRemoveCmd is the command to remove site authentication
type AuthRemoveCmd struct {
	Alias string `arg:"" help:"Site alias to remove"`
}

// Run executes the auth remove command
func (c *AuthRemoveCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// Check if site exists
	if _, ok := cfg.Sites[c.Alias]; !ok {
		return fmt.Errorf("site alias '%s' not found", c.Alias)
	}

	// Confirm (if --force flag is not set)
	if !root.Force {
		fmt.Printf("Remove site '%s'? (y/N): ", c.Alias)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// Remove from keyring
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	if err := store.Delete(c.Alias); err != nil {
		return err
	}

	// Remove from configuration
	delete(cfg.Sites, c.Alias)
	if cfg.DefaultSite == c.Alias {
		cfg.DefaultSite = ""
	}

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Removed site '%s'\n", c.Alias)
	return nil
}

// AuthStatusCmd is the command to check authentication status
type AuthStatusCmd struct{}

// Run executes the auth status command
func (c *AuthStatusCmd) Run(ctx context.Context, root *RootFlags) error {
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

	if len(cfg.Sites) == 0 {
		fmt.Println("No sites configured")
		fmt.Println("Run 'gho auth add <site-url>' to add a site")
		return nil
	}

	fmt.Printf("Default site: %s\n", cfg.DefaultSite)
	fmt.Printf("Configured sites: %d\n", len(cfg.Sites))
	fmt.Printf("Keyring backend: %s\n", cfg.KeyringBackend)

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "gho", "config.json"), nil
}

// getKeyringDir returns the path to the keyring directory
func getKeyringDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gho", "keyring")
}

// extractAliasFromURL extracts an alias from a URL
func extractAliasFromURL(url string) string {
	// https://myblog.ghost.io -> myblog
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	parts := strings.Split(url, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return "site"
}

// AuthTokensCmd is the root command for API token management
type AuthTokensCmd struct {
	List   AuthTokensListCmd   `cmd:"" help:"List API tokens"`
	Delete AuthTokensDeleteCmd `cmd:"" help:"Delete API token"`
}

// AuthTokensListCmd is the command to list API tokens
type AuthTokensListCmd struct{}

// Run executes the auth tokens list command
func (c *AuthTokensListCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// Open keyring store
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	// Get list of saved keys
	aliases, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list keys: %w", err)
	}

	if len(aliases) == 0 {
		fmt.Println("No API tokens found")
		return nil
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// For JSON format
	if root.JSON {
		type tokenInfo struct {
			Alias string `json:"alias"`
			KeyID string `json:"key_id"`
		}
		var tokens []tokenInfo

		for _, alias := range aliases {
			apiKey, err := store.Get(alias)
			if err != nil {
				continue
			}
			keyID, _, err := secrets.ParseAdminAPIKey(apiKey)
			if err != nil {
				continue
			}
			tokens = append(tokens, tokenInfo{
				Alias: alias,
				KeyID: keyID,
			})
		}

		return formatter.Print(tokens)
	}

	// Output in table format
	headers := []string{"Alias", "Key ID"}
	var rows [][]string

	for _, alias := range aliases {
		apiKey, err := store.Get(alias)
		if err != nil {
			continue
		}
		keyID, _, err := secrets.ParseAdminAPIKey(apiKey)
		if err != nil {
			continue
		}
		rows = append(rows, []string{alias, keyID})
	}

	return formatter.PrintTable(headers, rows)
}

// AuthTokensDeleteCmd is the command to delete API tokens
type AuthTokensDeleteCmd struct {
	Alias string `arg:"" help:"Site alias to delete token for"`
}

// Run executes the auth tokens delete command
func (c *AuthTokensDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
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

	// Open keyring store
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	// Remove from keyring
	if err := store.Delete(c.Alias); err != nil {
		return err
	}

	fmt.Printf("✓ Deleted API token for '%s' (configuration preserved)\n", c.Alias)
	return nil
}

// AuthCredentialsCmd is the command to add configuration from a credentials file
type AuthCredentialsCmd struct {
	File string `arg:"" help:"Credentials file path (JSON format)"`
}

// Run executes the auth credentials command
func (c *AuthCredentialsCmd) Run(ctx context.Context, root *RootFlags) error {
	// Load credentials file
	type credentialsFile struct {
		SiteURL string `json:"site_url"`
		Alias   string `json:"alias"`
		APIKey  string `json:"api_key"`
	}

	data, err := os.ReadFile(c.File)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds credentialsFile
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials file: %w", err)
	}

	// Check required fields
	if creds.SiteURL == "" {
		return fmt.Errorf("site_url is not specified")
	}
	if creds.APIKey == "" {
		return fmt.Errorf("api_key is not specified")
	}

	// Generate alias from URL if not provided
	alias := creds.Alias
	if alias == "" {
		alias = extractAliasFromURL(creds.SiteURL)
	}

	// Parse API key
	keyID, secret, err := secrets.ParseAdminAPIKey(creds.APIKey)
	if err != nil {
		return err
	}

	// Validate API key (by fetching site info)
	client, err := ghostapi.NewClient(creds.SiteURL, keyID, secret)
	if err != nil {
		return err
	}

	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}

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

	// Save to keyring
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	if err := store.Set(alias, creds.APIKey); err != nil {
		return err
	}

	// Add site to configuration
	cfg.AddSite(alias, creds.SiteURL)
	if cfg.DefaultSite == "" {
		cfg.DefaultSite = alias
	}

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Added site '%s' (%s)\n", alias, site.Title)
	if cfg.DefaultSite == alias {
		fmt.Println("✓ Set as default site")
	}

	return nil
}
