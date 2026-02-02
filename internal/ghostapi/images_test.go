/**
 * images_test.go
 * Test code for Images API
 */

package ghostapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUploadImage_UploadImage uploads an image
func TestUploadImage_UploadImage(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/images/upload/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/images/upload/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}

		// Verify Content-Type is multipart/form-data
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Errorf("Content-Type = %q; want multipart/form-data", contentType)
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			t.Fatalf("Failed to parse multipart form: %v", err)
		}

		// Verify file exists
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("Failed to retrieve file: %v", err)
		}
		defer file.Close()

		// Verify file name
		if header.Filename != "test-image.jpg" {
			t.Errorf("File name = %q; want %q", header.Filename, "test-image.jpg")
		}

		// Verify file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("Failed to read file content: %v", err)
		}
		expectedContent := "fake image content"
		if string(fileContent) != expectedContent {
			t.Errorf("File content = %q; want %q", string(fileContent), expectedContent)
		}

		// Return response
		response := map[string]interface{}{
			"images": []map[string]interface{}{
				{
					"url": "https://example.com/content/images/2024/01/test-image.jpg",
					"ref": nil,
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

	// Upload image
	fakeImageContent := strings.NewReader("fake image content")
	image, err := client.UploadImage(fakeImageContent, "test-image.jpg", ImageUploadOptions{})
	if err != nil {
		t.Fatalf("Failed to upload image: %v", err)
	}

	// Verify response
	expectedURL := "https://example.com/content/images/2024/01/test-image.jpg"
	if image.URL != expectedURL {
		t.Errorf("Image URL = %q; want %q", image.URL, expectedURL)
	}
}

// TestUploadImage_PurposeParameter tests purpose parameter
func TestUploadImage_PurposeParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			t.Fatalf("Failed to parse multipart form: %v", err)
		}

		// Verify purpose parameter
		purpose := r.FormValue("purpose")
		if purpose != "profile_image" {
			t.Errorf("purpose = %q; want %q", purpose, "profile_image")
		}

		// Verify ref parameter
		ref := r.FormValue("ref")
		if ref != "test-ref-12345" {
			t.Errorf("ref = %q; want %q", ref, "test-ref-12345")
		}

		// Return response
		response := map[string]interface{}{
			"images": []map[string]interface{}{
				{
					"url": "https://example.com/content/images/2024/01/profile.jpg",
					"ref": "test-ref-12345",
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

	// Upload image with purpose and ref
	fakeImageContent := strings.NewReader("fake profile image")
	image, err := client.UploadImage(fakeImageContent, "profile.jpg", ImageUploadOptions{
		Purpose: "profile_image",
		Ref:     "test-ref-12345",
	})
	if err != nil {
		t.Fatalf("Failed to upload image: %v", err)
	}

	// Verify response
	if image.Ref != "test-ref-12345" {
		t.Errorf("ref = %q; want %q", image.Ref, "test-ref-12345")
	}
}
