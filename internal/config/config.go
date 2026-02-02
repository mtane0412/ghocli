/**
 * config.go
 * Configuration file management for gho
 *
 * The configuration file is saved at ~/.config/gho/config.json
 * and provides alias functionality for multi-site support.
 */

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the gho configuration
type Config struct {
	// KeyringBackend is the keyring backend type (auto/file/keychain, etc.)
	KeyringBackend string `json:"keyring_backend"`

	// DefaultSite is the default site alias
	DefaultSite string `json:"default_site,omitempty"`

	// Sites is a mapping from alias to site URL
	Sites map[string]string `json:"sites"`
}

// Load reads the configuration file from the specified path.
// If the file does not exist, it returns a new configuration with default values.
func Load(path string) (*Config, error) {
	// If file does not exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			KeyringBackend: "auto",
			DefaultSite:    "",
			Sites:          make(map[string]string),
		}, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Initialize Sites if nil
	if cfg.Sites == nil {
		cfg.Sites = make(map[string]string)
	}

	return &cfg, nil
}

// Save writes the configuration to the specified path.
func (c *Config) Save(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Convert to JSON (with indentation)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to file (0600 = read/write for owner only)
	return os.WriteFile(path, data, 0600)
}

// AddSite adds a site alias and URL to the configuration.
func (c *Config) AddSite(alias, url string) {
	if c.Sites == nil {
		c.Sites = make(map[string]string)
	}
	c.Sites[alias] = url
}

// GetSiteURL retrieves a site URL from an alias or URL string.
// If registered as an alias, it returns the corresponding URL.
// Otherwise, it treats it as a URL string and returns it as is.
func (c *Config) GetSiteURL(aliasOrURL string) (string, bool) {
	// Check if registered as an alias
	if url, ok := c.Sites[aliasOrURL]; ok {
		return url, true
	}

	// Treat as URL string (if it starts with https://)
	if strings.HasPrefix(aliasOrURL, "https://") || strings.HasPrefix(aliasOrURL, "http://") {
		return aliasOrURL, true
	}

	return "", false
}
