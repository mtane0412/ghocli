/**
 * store.go
 * Secure storage of Admin API keys via keyring integration
 *
 * Uses 99designs/keyring to store API keys in OS keyring.
 * Supports macOS: Keychain, Linux: Secret Service, Windows: Credential Manager.
 */

package secrets

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

// Store is a store for saving and retrieving API keys
type Store struct {
	ring keyring.Keyring
}

const (
	// ServiceName is the service name used when saving to keyring
	ServiceName = "gho-ghost-admin"
)

// NewStore creates a new keyring store.
// backend: "auto", "file", "keychain", etc.
// fileDir: directory for file storage when backend is "file"
func NewStore(backend, fileDir string) (*Store, error) {
	var cfg keyring.Config
	cfg.ServiceName = ServiceName

	// Override backend with environment variable
	if envBackend := os.Getenv("GHO_KEYRING_BACKEND"); envBackend != "" {
		backend = envBackend
	}

	// Set backend type
	switch backend {
	case "auto":
		cfg.AllowedBackends = []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		}
	case "file":
		cfg.AllowedBackends = []keyring.BackendType{keyring.FileBackend}
		cfg.FileDir = fileDir
		cfg.FilePasswordFunc = func(prompt string) (string, error) {
			// Get password from environment variable
			if pw := os.Getenv("GHO_KEYRING_PASSWORD"); pw != "" {
				return pw, nil
			}
			// Return empty string if password is not set
			return "", nil
		}
	case "keychain":
		cfg.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
	case "secretservice":
		cfg.AllowedBackends = []keyring.BackendType{keyring.SecretServiceBackend}
	case "wincred":
		cfg.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
	default:
		return nil, fmt.Errorf("unknown backend: %s", backend)
	}

	ring, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &Store{ring: ring}, nil
}

// Set saves an Admin API key for a site alias.
func (s *Store) Set(alias, apiKey string) error {
	item := keyring.Item{
		Key:  alias,
		Data: []byte(apiKey),
	}

	if err := s.ring.Set(item); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}

	return nil
}

// Get retrieves the Admin API key for a site alias.
func (s *Store) Get(alias string) (string, error) {
	item, err := s.ring.Get(alias)
	if err != nil {
		return "", fmt.Errorf("failed to get API key: %w", err)
	}

	return string(item.Data), nil
}

// Delete removes the Admin API key for a site alias.
func (s *Store) Delete(alias string) error {
	if err := s.ring.Remove(alias); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

// List retrieves all site aliases stored in the keyring.
func (s *Store) List() ([]string, error) {
	keys, err := s.ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	return keys, nil
}

// ParseAdminAPIKey parses a Ghost Admin API key (id:secret format).
func ParseAdminAPIKey(apiKey string) (id, secret string, err error) {
	if apiKey == "" {
		return "", "", errors.New("API key is empty")
	}

	parts := strings.SplitN(apiKey, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid API key format (must be id:secret)")
	}

	return parts[0], parts[1], nil
}
