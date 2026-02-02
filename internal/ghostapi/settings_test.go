/**
 * settings_test.go
 * Test code for Settings API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetSettings_RetrieveSettingsList retrieves the list of settings
func TestGetSettings_RetrieveSettingsList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/settings/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/settings/")
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Verify Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}

		// Return response
		response := map[string]interface{}{
			"settings": []map[string]interface{}{
				{
					"key":   "title",
					"value": "My Ghost Site",
				},
				{
					"key":   "description",
					"value": "Thoughts, stories and ideas",
				},
				{
					"key":   "timezone",
					"value": "Asia/Tokyo",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Retrieve settings list
	resp, err := client.GetSettings()
	if err != nil {
		t.Fatalf("Failed to retrieve settings list: %v", err)
	}

	// Verify response
	if len(resp.Settings) != 3 {
		t.Errorf("Number of settings = %d; want 3", len(resp.Settings))
	}

	// Verify first setting
	firstSetting := resp.Settings[0]
	if firstSetting.Key != "title" {
		t.Errorf("Setting key = %q; want %q", firstSetting.Key, "title")
	}
	if firstSetting.Value != "My Ghost Site" {
		t.Errorf("Setting value = %q; want %q", firstSetting.Value, "My Ghost Site")
	}
}

// TestUpdateSettings_UpdateSettings updates settings
func TestUpdateSettings_UpdateSettings(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/settings/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/settings/")
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Parse request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify settings field exists
		settings, ok := reqBody["settings"].([]interface{})
		if !ok {
			t.Fatal("settings field does not exist")
		}

		if len(settings) != 1 {
			t.Errorf("Number of settings = %d; want 1", len(settings))
		}

		// Return response
		response := map[string]interface{}{
			"settings": []map[string]interface{}{
				{
					"key":   "title",
					"value": "Updated Title",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Update settings
	updates := []SettingUpdate{
		{
			Key:   "title",
			Value: "Updated Title",
		},
	}
	resp, err := client.UpdateSettings(updates)
	if err != nil {
		t.Fatalf("Failed to update settings: %v", err)
	}

	// Verify response
	if len(resp.Settings) != 1 {
		t.Errorf("Number of settings = %d; want 1", len(resp.Settings))
	}

	updatedSetting := resp.Settings[0]
	if updatedSetting.Key != "title" {
		t.Errorf("Setting key = %q; want %q", updatedSetting.Key, "title")
	}
	if updatedSetting.Value != "Updated Title" {
		t.Errorf("Setting value = %q; want %q", updatedSetting.Value, "Updated Title")
	}
}
