/**
 * themes_test.go
 * Test code for Themes API
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListThemes_RetrieveThemeList retrieves the list of themes
func TestListThemes_RetrieveThemeList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/themes/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/themes/")
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
			"themes": []map[string]interface{}{
				{
					"name":   "casper",
					"active": true,
					"package": map[string]interface{}{
						"name":        "casper",
						"description": "The default theme for Ghost",
						"version":     "5.0.0",
					},
					"templates": []map[string]interface{}{
						{"filename": "index.hbs"},
						{"filename": "post.hbs"},
					},
				},
				{
					"name":   "starter",
					"active": false,
					"package": map[string]interface{}{
						"name":        "starter",
						"description": "A minimal starter theme",
						"version":     "1.0.0",
					},
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

	// Retrieve theme list
	resp, err := client.ListThemes()
	if err != nil {
		t.Fatalf("Failed to retrieve theme list: %v", err)
	}

	// Verify response
	if len(resp.Themes) != 2 {
		t.Errorf("Number of themes = %d; want 2", len(resp.Themes))
	}

	// Verify first theme
	firstTheme := resp.Themes[0]
	if firstTheme.Name != "casper" {
		t.Errorf("Theme name = %q; want %q", firstTheme.Name, "casper")
	}
	if !firstTheme.Active {
		t.Error("Active flag = false; want true")
	}
	if firstTheme.Package == nil {
		t.Fatal("Package information is nil")
	}
	if firstTheme.Package.Name != "casper" {
		t.Errorf("Package name = %q; want %q", firstTheme.Package.Name, "casper")
	}
	if len(firstTheme.Templates) != 2 {
		t.Errorf("Number of templates = %d; want 2", len(firstTheme.Templates))
	}
}

// TestUploadTheme_UploadTheme uploads a theme
func TestUploadTheme_UploadTheme(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/themes/upload/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/themes/upload/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify Content-Type is multipart/form-data
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || len(contentType) < 19 || contentType[:19] != "multipart/form-data" {
			t.Errorf("Content-Type = %q; want multipart/form-data", contentType)
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("Failed to parse multipart form: %v", err)
		}

		// Verify file exists
		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("Failed to retrieve file: %v", err)
		}
		if fileHeader.Filename != "theme.zip" {
			t.Errorf("File name = %q; want %q", fileHeader.Filename, "theme.zip")
		}

		// Return response
		response := map[string]interface{}{
			"themes": []map[string]interface{}{
				{
					"name":   "custom-theme",
					"active": false,
					"package": map[string]interface{}{
						"name":        "custom-theme",
						"description": "A custom theme",
						"version":     "1.0.0",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create dummy ZIP file data
	fileData := []byte("dummy zip content")
	reader := bytes.NewReader(fileData)

	// Upload theme
	theme, err := client.UploadTheme(reader, "theme.zip")
	if err != nil {
		t.Fatalf("Failed to upload theme: %v", err)
	}

	// Verify response
	if theme.Name != "custom-theme" {
		t.Errorf("Theme name = %q; want %q", theme.Name, "custom-theme")
	}
	if theme.Active {
		t.Error("Active flag = true; want false")
	}
	if theme.Package == nil {
		t.Fatal("Package information is nil")
	}
	if theme.Package.Version != "1.0.0" {
		t.Errorf("Version = %q; want %q", theme.Package.Version, "1.0.0")
	}
}

// TestActivateTheme_ActivateTheme activates a theme
func TestActivateTheme_ActivateTheme(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/themes/custom-theme/activate/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Return response
		response := map[string]interface{}{
			"themes": []map[string]interface{}{
				{
					"name":   "custom-theme",
					"active": true,
					"package": map[string]interface{}{
						"name":        "custom-theme",
						"description": "A custom theme",
						"version":     "1.0.0",
					},
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

	// Activate theme
	theme, err := client.ActivateTheme("custom-theme")
	if err != nil {
		t.Fatalf("Failed to activate theme: %v", err)
	}

	// Verify response
	if theme.Name != "custom-theme" {
		t.Errorf("Theme name = %q; want %q", theme.Name, "custom-theme")
	}
	if !theme.Active {
		t.Error("Active flag = false; want true")
	}
}

// TestDeleteTheme_DeleteTheme deletes a theme
func TestDeleteTheme_DeleteTheme(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/themes/custom-theme/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "DELETE")
		}

		// Return response (204 No Content)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Delete theme
	err = client.DeleteTheme("custom-theme")
	if err != nil {
		t.Fatalf("Failed to delete theme: %v", err)
	}
}
