/**
 * tags_test.go
 * Test code for Tags API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListTags_GetTagList tests retrieving a list of tags
func TestListTags_GetTagList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/tags/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/tags/")
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
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "Technology related articles",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
				},
				{
					"id":          "64fac5417c4c6b0001234568",
					"name":        "Programming",
					"slug":        "programming",
					"description": "Programming tips",
					"visibility":  "public",
					"created_at":  "2024-01-16T10:00:00.000Z",
					"updated_at":  "2024-01-16T10:00:00.000Z",
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 2,
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

	// Get tag list
	resp, err := client.ListTags(TagListOptions{})
	if err != nil {
		t.Fatalf("Failed to get tag list: %v", err)
	}

	// Verify response
	if len(resp.Tags) != 2 {
		t.Errorf("Number of tags = %d; want 2", len(resp.Tags))
	}

	// Verify first tag
	firstTag := resp.Tags[0]
	if firstTag.Name != "Technology" {
		t.Errorf("Tag name = %q; want %q", firstTag.Name, "Technology")
	}
	if firstTag.Slug != "technology" {
		t.Errorf("Slug = %q; want %q", firstTag.Slug, "technology")
	}
}

// TestListTags_IncludeParameter tests the include parameter
func TestListTags_IncludeParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if !r.URL.Query().Has("include") {
			t.Error("include parameter is not set")
		}
		if r.URL.Query().Get("include") != "count.posts" {
			t.Errorf("include parameter = %q; want %q", r.URL.Query().Get("include"), "count.posts")
		}

		// Return response
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
					"count": map[string]interface{}{
						"posts": 10,
					},
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 1,
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

	// Get tag list (including count.posts)
	resp, err := client.ListTags(TagListOptions{
		Include: "count.posts",
	})
	if err != nil {
		t.Fatalf("Failed to get tag list: %v", err)
	}

	// Verify response
	if len(resp.Tags) != 1 {
		t.Errorf("Number of tags = %d; want 1", len(resp.Tags))
	}
}

// TestGetTag_GetByID tests retrieving a tag by ID
func TestGetTag_GetByID(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "Technology related articles",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
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

	// Get tag
	tag, err := client.GetTag("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Failed to get tag: %v", err)
	}

	// Verify response
	if tag.Name != "Technology" {
		t.Errorf("Tag name = %q; want %q", tag.Name, "Technology")
	}
	if tag.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("Tag ID = %q; want %q", tag.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetTag_GetBySlug tests retrieving a tag by slug
func TestGetTag_GetBySlug(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/tags/slug/technology/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "Technology related articles",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
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

	// Get tag
	tag, err := client.GetTag("slug:technology")
	if err != nil {
		t.Fatalf("Failed to get tag: %v", err)
	}

	// Verify response
	if tag.Slug != "technology" {
		t.Errorf("Slug = %q; want %q", tag.Slug, "technology")
	}
}

// TestCreateTag_CreateTag tests creating a tag
func TestCreateTag_CreateTag(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/tags/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/tags/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		tags, ok := reqBody["tags"].([]interface{})
		if !ok || len(tags) == 0 {
			t.Error("Request body does not contain tags array")
		}

		// Return response
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234999",
					"name":        "New Tag",
					"slug":        "new-tag",
					"description": "A new tag",
					"visibility":  "public",
					"created_at":  createdAt,
					"updated_at":  createdAt,
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

	// Create tag
	newTag := &Tag{
		Name:        "New Tag",
		Description: "A new tag",
	}

	createdTag, err := client.CreateTag(newTag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	// Verify response
	if createdTag.Name != "New Tag" {
		t.Errorf("Tag name = %q; want %q", createdTag.Name, "New Tag")
	}
	if createdTag.ID == "" {
		t.Error("Tag ID is empty")
	}
}

// TestUpdateTag_UpdateTag tests updating a tag
func TestUpdateTag_UpdateTag(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Return response
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Updated Technology",
					"slug":        "technology",
					"description": "Updated technology related articles",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  time.Now().Format(time.RFC3339),
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

	// Update tag
	updateTag := &Tag{
		Name:        "Updated Technology",
		Description: "Updated technology related articles",
	}

	updatedTag, err := client.UpdateTag("64fac5417c4c6b0001234567", updateTag)
	if err != nil {
		t.Fatalf("Failed to update tag: %v", err)
	}

	// Verify response
	if updatedTag.Name != "Updated Technology" {
		t.Errorf("Tag name = %q; want %q", updatedTag.Name, "Updated Technology")
	}
}

// TestDeleteTag_DeleteTag tests deleting a tag
func TestDeleteTag_DeleteTag(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "DELETE")
		}

		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Delete tag
	err = client.DeleteTag("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Failed to delete tag: %v", err)
	}
}
