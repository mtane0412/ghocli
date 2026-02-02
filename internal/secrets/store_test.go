/**
 * store_test.go
 * Test code for keyring integration
 */

package secrets

import (
	"os"
	"testing"
)

// TestStore_BasicSetAndGet tests basic Set and Get operations
func TestStore_BasicSetAndGet(t *testing.T) {
	// Create test store (using file backend)
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Save API key
	testKey := "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	if err := store.Set("testsite", testKey); err != nil {
		t.Fatalf("failed to save API key: %v", err)
	}

	// Retrieve API key
	retrieved, err := store.Get("testsite")
	if err != nil {
		t.Fatalf("failed to retrieve API key: %v", err)
	}

	// Verify saved and retrieved keys match
	if retrieved != testKey {
		t.Errorf("retrieved key = %q; want %q", retrieved, testKey)
	}
}

// TestStore_GetNonexistentKey tests retrieving a nonexistent key
func TestStore_GetNonexistentKey(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Attempt to retrieve nonexistent key
	_, err = store.Get("nonexistent")
	if err == nil {
		t.Error("expected error when retrieving nonexistent key")
	}
}

// TestStore_DeleteKey tests key deletion
func TestStore_DeleteKey(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Save API key
	testKey := "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	if err := store.Set("testsite", testKey); err != nil {
		t.Fatalf("failed to save API key: %v", err)
	}

	// Delete key
	if err := store.Delete("testsite"); err != nil {
		t.Fatalf("failed to delete API key: %v", err)
	}

	// Verify key cannot be retrieved after deletion
	_, err = store.Get("testsite")
	if err == nil {
		t.Error("deleted key should not be retrievable")
	}
}

// TestStore_ListSavedKeys tests listing all saved keys
func TestStore_ListSavedKeys(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Save multiple API keys
	keys := map[string]string{
		"site1": "key1:secret1",
		"site2": "key2:secret2",
		"site3": "key3:secret3",
	}
	for alias, key := range keys {
		if err := store.Set(alias, key); err != nil {
			t.Fatalf("failed to save API key (%s): %v", alias, err)
		}
	}

	// Retrieve list of saved keys
	aliases, err := store.List()
	if err != nil {
		t.Fatalf("failed to retrieve key list: %v", err)
	}

	// Verify all aliases are included
	if len(aliases) != len(keys) {
		t.Errorf("key count = %d; want %d", len(aliases), len(keys))
	}

	for alias := range keys {
		found := false
		for _, a := range aliases {
			if a == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("alias %q not found in list", alias)
		}
	}
}

// TestStore_ParseAdminAPIKey tests parsing Admin API keys
func TestStore_ParseAdminAPIKey(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		wantID    string
		wantSecret string
		wantErr   bool
	}{
		{
			name:       "valid format",
			input:      "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123",
			wantID:     "64fac5417c4c6b0001234567",
			wantSecret: "89abcdef01234567890123456789abcd01234567890123456789abcdef0123",
			wantErr:    false,
		},
		{
			name:    "missing colon",
			input:   "64fac5417c4c6b000123456789abcdef",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, secret, err := ParseAdminAPIKey(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tc.wantID {
				t.Errorf("id = %q; want %q", id, tc.wantID)
			}
			if secret != tc.wantSecret {
				t.Errorf("secret = %q; want %q", secret, tc.wantSecret)
			}
		})
	}
}

// TestNewStore_GHO_KEYRING_BACKEND tests that GHO_KEYRING_BACKEND environment variable overrides backend
func TestNewStore_GHO_KEYRING_BACKEND(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_KEYRING_BACKEND", "file")
	defer os.Unsetenv("GHO_KEYRING_BACKEND")

	// Even with "auto" as backend argument, environment variable should force "file" backend
	store, err := NewStore("auto", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Verify store is properly created (basic operations work)
	testKey := "test:key"
	if err := store.Set("test", testKey); err != nil {
		t.Fatalf("failed to save API key: %v", err)
	}

	retrieved, err := store.Get("test")
	if err != nil {
		t.Fatalf("failed to retrieve API key: %v", err)
	}

	if retrieved != testKey {
		t.Errorf("retrieved key = %q; want %q", retrieved, testKey)
	}
}

// TestNewStore_GHO_KEYRING_PASSWORD tests that GHO_KEYRING_PASSWORD environment variable provides password
func TestNewStore_GHO_KEYRING_PASSWORD(t *testing.T) {
	// Set password environment variable
	testPassword := "test-password"
	os.Setenv("GHO_KEYRING_PASSWORD", testPassword)
	defer os.Unsetenv("GHO_KEYRING_PASSWORD")

	// Create store using file backend
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Verify store is properly created (basic operations work)
	testKey := "test:key"
	if err := store.Set("test", testKey); err != nil {
		t.Fatalf("failed to save API key: %v", err)
	}

	retrieved, err := store.Get("test")
	if err != nil {
		t.Fatalf("failed to retrieve API key: %v", err)
	}

	if retrieved != testKey {
		t.Errorf("retrieved key = %q; want %q", retrieved, testKey)
	}
}
