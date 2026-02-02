/**
 * newsletters_test.go
 * Test code for Newsletters API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListNewsletters_GetNewsletterList retrieves a list of newsletters
func TestListNewsletters_GetNewsletterList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/ghost/api/admin/newsletters/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/newsletters/")
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Verify that Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "Main Newsletter",
					"description":         "Weekly Newsletter",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost Editorial Team",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
				},
				{
					"id":                  "64fac5417c4c6b0001234568",
					"name":                "Premium Newsletter",
					"description":         "Paid Members Only Newsletter",
					"slug":                "premium-newsletter",
					"status":              "active",
					"visibility":          "paid",
					"subscribe_on_signup": false,
					"sender_name":         "Ghost Editorial Team",
					"sender_email":        "premium@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          1,
					"created_at":          "2024-01-16T10:00:00.000Z",
					"updated_at":          "2024-01-16T10:00:00.000Z",
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

	// Get newsletter list
	resp, err := client.ListNewsletters(NewsletterListOptions{})
	if err != nil {
		t.Fatalf("Failed to get newsletter list: %v", err)
	}

	// Verify the response
	if len(resp.Newsletters) != 2 {
		t.Errorf("Number of newsletters = %d; want 2", len(resp.Newsletters))
	}

	// Verify the first newsletter
	firstNewsletter := resp.Newsletters[0]
	if firstNewsletter.Name != "Main Newsletter" {
		t.Errorf("Newsletter name = %q; want %q", firstNewsletter.Name, "Main Newsletter")
	}
	if firstNewsletter.Slug != "main-newsletter" {
		t.Errorf("Slug = %q; want %q", firstNewsletter.Slug, "main-newsletter")
	}
	if firstNewsletter.Status != "active" {
		t.Errorf("Status = %q; want %q", firstNewsletter.Status, "active")
	}
}

// TestListNewsletters_FilterParameter tests filter parameter usage
func TestListNewsletters_FilterParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if !r.URL.Query().Has("filter") {
			t.Error("filter parameter is not set")
		}
		if r.URL.Query().Get("filter") != "status:active" {
			t.Errorf("filter parameter = %q; want %q", r.URL.Query().Get("filter"), "status:active")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "Main Newsletter",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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

	// Get newsletter list (filtered by status:active)
	resp, err := client.ListNewsletters(NewsletterListOptions{
		Filter: "status:active",
	})
	if err != nil {
		t.Fatalf("Failed to get newsletter list: %v", err)
	}

	// Verify the response
	if len(resp.Newsletters) != 1 {
		t.Errorf("Number of newsletters = %d; want 1", len(resp.Newsletters))
	}
}

// TestGetNewsletter_GetNewsletterByID retrieves a newsletter by ID
func TestGetNewsletter_GetNewsletterByID(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/newsletters/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "Main Newsletter",
					"description":         "Weekly Newsletter",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost Editorial Team",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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

	// Get newsletter
	newsletter, err := client.GetNewsletter("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Failed to get newsletter: %v", err)
	}

	// Verify the response
	if newsletter.Name != "Main Newsletter" {
		t.Errorf("Newsletter name = %q; want %q", newsletter.Name, "Main Newsletter")
	}
	if newsletter.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("Newsletter ID = %q; want %q", newsletter.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetNewsletter_GetNewsletterBySlug retrieves a newsletter by slug
func TestGetNewsletter_GetNewsletterBySlug(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/newsletters/slug/main-newsletter/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "Main Newsletter",
					"description":         "Weekly Newsletter",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost Editorial Team",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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

	// Get newsletter
	newsletter, err := client.GetNewsletter("slug:main-newsletter")
	if err != nil {
		t.Fatalf("Failed to get newsletter: %v", err)
	}

	// Verify the response
	if newsletter.Slug != "main-newsletter" {
		t.Errorf("Slug = %q; want %q", newsletter.Slug, "main-newsletter")
	}
}

// TestCreateNewsletter_CreateNewsletter creates a new newsletter
func TestCreateNewsletter_CreateNewsletter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/ghost/api/admin/newsletters/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/newsletters/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		newsletters, ok := reqBody["newsletters"].([]interface{})
		if !ok || len(newsletters) == 0 {
			t.Error("newsletters field is invalid")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234569",
					"name":                "New Newsletter",
					"description":         "Test Newsletter",
					"slug":                "new-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Test Editorial Team",
					"sender_email":        "test@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-20T10:00:00.000Z",
					"updated_at":          "2024-01-20T10:00:00.000Z",
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

	// Create newsletter
	newNewsletter := &Newsletter{
		Name:              "New Newsletter",
		Description:       "Test Newsletter",
		Visibility:        "members",
		SubscribeOnSignup: true,
		SenderName:        "Test Editorial Team",
		SenderEmail:       "test@example.com",
	}

	createdNewsletter, err := client.CreateNewsletter(newNewsletter)
	if err != nil {
		t.Fatalf("Failed to create newsletter: %v", err)
	}

	// Verify the response
	if createdNewsletter.Name != "New Newsletter" {
		t.Errorf("Newsletter name = %q; want %q", createdNewsletter.Name, "New Newsletter")
	}
	if createdNewsletter.ID == "" {
		t.Error("Newsletter ID is not set")
	}
}

// TestUpdateNewsletter_UpdateNewsletter updates an existing newsletter
func TestUpdateNewsletter_UpdateNewsletter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/newsletters/64fac5417c4c6b0001234567/"
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

		newsletters, ok := reqBody["newsletters"].([]interface{})
		if !ok || len(newsletters) == 0 {
			t.Error("newsletters field is invalid")
		}

		// Return response
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "Updated Newsletter",
					"description":         "Updated Description",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Updated Editorial Team",
					"sender_email":        "updated@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-20T10:00:00.000Z",
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

	// Update newsletter
	updateNewsletter := &Newsletter{
		Name:        "Updated Newsletter",
		Description: "Updated Description",
		SenderName:  "Updated Editorial Team",
		SenderEmail: "updated@example.com",
	}

	updatedNewsletter, err := client.UpdateNewsletter("64fac5417c4c6b0001234567", updateNewsletter)
	if err != nil {
		t.Fatalf("Failed to update newsletter: %v", err)
	}

	// Verify the response
	if updatedNewsletter.Name != "Updated Newsletter" {
		t.Errorf("Newsletter name = %q; want %q", updatedNewsletter.Name, "Updated Newsletter")
	}
	if updatedNewsletter.Description != "Updated Description" {
		t.Errorf("Description = %q; want %q", updatedNewsletter.Description, "Updated Description")
	}
}
