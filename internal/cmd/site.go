/**
 * site.go
 * Site information retrieval command
 *
 * Retrieves and displays basic information about a Ghost site.
 */

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/config"
	"github.com/mtane0412/ghocli/internal/errfmt"
	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
	"github.com/mtane0412/ghocli/internal/secrets"
)

// SiteCmd is the command to retrieve site information
type SiteCmd struct{}

// Run executes the site command
func (c *SiteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get site information
	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("failed to get site information: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(site)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"title", site.Title},
		{"description", site.Description},
		{"url", site.URL},
		{"version", site.Version},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// getAPIClient retrieves an API client
func getAPIClient(root *RootFlags) (*ghostapi.Client, error) {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Load config
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Determine site URL
	siteURL := root.Site
	if siteURL == "" {
		siteURL = cfg.DefaultSite
	}
	if siteURL == "" {
		return nil, errors.New(errfmt.FormatSiteError())
	}

	// Convert alias to URL if applicable
	if url, ok := cfg.GetSiteURL(siteURL); ok {
		siteURL = url
	} else {
		return nil, fmt.Errorf("site '%s' not found", siteURL)
	}

	// Reverse lookup alias
	alias := ""
	for a, u := range cfg.Sites {
		if u == siteURL {
			alias = a
			break
		}
	}
	if alias == "" {
		return nil, fmt.Errorf("alias not found for site URL")
	}

	// Get API key from keyring
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	apiKey, err := store.Get(alias)
	if err != nil {
		return nil, errors.New(errfmt.FormatAuthError(alias))
	}

	// Parse API key
	keyID, secret, err := secrets.ParseAdminAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	// Create API client
	return ghostapi.NewClient(siteURL, keyID, secret)
}
