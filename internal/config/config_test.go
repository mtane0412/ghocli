/**
 * config_test.go
 * Test code for configuration system
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig_DefaultValuesOnNewCreation tests default values when creating a new config file
func TestLoadConfig_DefaultValuesOnNewCreation(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Call Load when the config file does not exist
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	}

	// Verify default values
	if cfg.KeyringBackend != "auto" {
		t.Errorf("KeyringBackend = %q; want %q", cfg.KeyringBackend, "auto")
	}
	if cfg.DefaultSite != "" {
		t.Errorf("DefaultSite = %q; want empty string", cfg.DefaultSite)
	}
	if cfg.Sites == nil {
		t.Error("Sites map is nil; want empty map")
	}
}

// TestLoadConfig_LoadExistingFile tests loading an existing config file
func TestLoadConfig_LoadExistingFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create a test config file
	configContent := `{
  "keyring_backend": "file",
  "default_site": "myblog",
  "sites": {
    "myblog": "https://myblog.ghost.io",
    "company": "https://blog.company.com"
  }
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load the config file
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	}

	// Verify loaded values
	if cfg.KeyringBackend != "file" {
		t.Errorf("KeyringBackend = %q; want %q", cfg.KeyringBackend, "file")
	}
	if cfg.DefaultSite != "myblog" {
		t.Errorf("DefaultSite = %q; want %q", cfg.DefaultSite, "myblog")
	}
	if len(cfg.Sites) != 2 {
		t.Errorf("Sites has %d entries; want 2", len(cfg.Sites))
	}
	if cfg.Sites["myblog"] != "https://myblog.ghost.io" {
		t.Errorf("Sites[myblog] = %q; want %q", cfg.Sites["myblog"], "https://myblog.ghost.io")
	}
	if cfg.Sites["company"] != "https://blog.company.com" {
		t.Errorf("Sites[company] = %q; want %q", cfg.Sites["company"], "https://blog.company.com")
	}
}

// TestSave_SaveConfigFile tests saving a config file
func TestSave_SaveConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create a config
	cfg := &Config{
		KeyringBackend: "auto",
		DefaultSite:    "testsite",
		Sites: map[string]string{
			"testsite": "https://test.ghost.io",
		},
	}

	// Save the config
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Failed to save config file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Reload the saved config and verify
	reloaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload saved config file: %v", err)
	}

	if reloaded.KeyringBackend != cfg.KeyringBackend {
		t.Errorf("KeyringBackend = %q; want %q", reloaded.KeyringBackend, cfg.KeyringBackend)
	}
	if reloaded.DefaultSite != cfg.DefaultSite {
		t.Errorf("DefaultSite = %q; want %q", reloaded.DefaultSite, cfg.DefaultSite)
	}
	if reloaded.Sites["testsite"] != cfg.Sites["testsite"] {
		t.Errorf("Sites[testsite] = %q; want %q", reloaded.Sites["testsite"], cfg.Sites["testsite"])
	}
}

// TestAddSite_AddSite tests adding a site
func TestAddSite_AddSite(t *testing.T) {
	cfg := &Config{
		KeyringBackend: "auto",
		Sites:          make(map[string]string),
	}

	// Add a site
	cfg.AddSite("myblog", "https://myblog.ghost.io")

	// Verify the site was added
	if cfg.Sites["myblog"] != "https://myblog.ghost.io" {
		t.Errorf("Sites[myblog] = %q; want %q", cfg.Sites["myblog"], "https://myblog.ghost.io")
	}
}

// TestGetSiteURL_GetURLFromAlias tests getting URL from an alias
func TestGetSiteURL_GetURLFromAlias(t *testing.T) {
	cfg := &Config{
		Sites: map[string]string{
			"myblog": "https://myblog.ghost.io",
		},
	}

	// Get URL from alias
	url, ok := cfg.GetSiteURL("myblog")
	if !ok {
		t.Fatal("GetSiteURL returned false; want true")
	}
	if url != "https://myblog.ghost.io" {
		t.Errorf("url = %q; want %q", url, "https://myblog.ghost.io")
	}

	// Non-existent alias
	_, ok = cfg.GetSiteURL("nonexistent")
	if ok {
		t.Error("GetSiteURL returned true for nonexistent alias; want false")
	}
}

// TestGetSiteURL_DirectURLSpecification tests direct URL specification
func TestGetSiteURL_DirectURLSpecification(t *testing.T) {
	cfg := &Config{
		Sites: make(map[string]string),
	}

	// Direct URL specification (returns as-is if not registered as an alias)
	url, ok := cfg.GetSiteURL("https://direct.ghost.io")
	if !ok {
		t.Fatal("GetSiteURL returned false for direct URL; want true")
	}
	if url != "https://direct.ghost.io" {
		t.Errorf("url = %q; want %q", url, "https://direct.ghost.io")
	}
}
