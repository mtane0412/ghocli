/**
 * webhooks_test.go
 * Test code for Webhooks API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreateWebhook_CreateWebhook tests the creation of a webhook
func TestCreateWebhook_CreateWebhook(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/webhooks/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/webhooks/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify that Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("request body parse error: %v", err)
		}

		webhooks, ok := reqBody["webhooks"].([]interface{})
		if !ok || len(webhooks) == 0 {
			t.Error("webhooks array does not exist in request body")
		}

		// Return response
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"webhooks": []map[string]interface{}{
				{
					"id":                "64fac5417c4c6b0001234567",
					"event":             "post.published",
					"target_url":        "https://example.com/webhook",
					"name":              "Post published webhook",
					"secret":            "secret123",
					"api_version":       "v5.0",
					"integration_id":    "64fac5417c4c6b0001234568",
					"status":            "available",
					"last_triggered_at": nil,
					"created_at":        createdAt,
					"updated_at":        createdAt,
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
		t.Fatalf("client creation error: %v", err)
	}

	// Create webhook
	webhook := &Webhook{
		Event:     "post.published",
		TargetURL: "https://example.com/webhook",
		Name:      "Post published webhook",
	}

	created, err := client.CreateWebhook(webhook)
	if err != nil {
		t.Fatalf("webhook creation error: %v", err)
	}

	// Verify response
	if created.Event != "post.published" {
		t.Errorf("event = %q; want %q", created.Event, "post.published")
	}
	if created.TargetURL != "https://example.com/webhook" {
		t.Errorf("target URL = %q; want %q", created.TargetURL, "https://example.com/webhook")
	}
	if created.ID == "" {
		t.Error("webhook ID is empty")
	}
	if created.Status != "available" {
		t.Errorf("status = %q; want %q", created.Status, "available")
	}
}

// TestUpdateWebhook_UpdateWebhook tests the update of a webhook
func TestUpdateWebhook_UpdateWebhook(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/webhooks/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("request body parse error: %v", err)
		}

		// Return response
		updatedAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"webhooks": []map[string]interface{}{
				{
					"id":                "64fac5417c4c6b0001234567",
					"event":             "post.published",
					"target_url":        "https://example.com/webhook-updated",
					"name":              "Updated webhook",
					"secret":            "secret123",
					"api_version":       "v5.0",
					"integration_id":    "64fac5417c4c6b0001234568",
					"status":            "available",
					"last_triggered_at": nil,
					"created_at":        "2024-01-15T10:00:00.000Z",
					"updated_at":        updatedAt,
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
		t.Fatalf("client creation error: %v", err)
	}

	// Update webhook
	updateWebhook := &Webhook{
		TargetURL: "https://example.com/webhook-updated",
		Name:      "Updated webhook",
	}

	updated, err := client.UpdateWebhook("64fac5417c4c6b0001234567", updateWebhook)
	if err != nil {
		t.Fatalf("webhook update error: %v", err)
	}

	// Verify response
	if updated.TargetURL != "https://example.com/webhook-updated" {
		t.Errorf("target URL = %q; want %q", updated.TargetURL, "https://example.com/webhook-updated")
	}
	if updated.Name != "Updated webhook" {
		t.Errorf("name = %q; want %q", updated.Name, "Updated webhook")
	}
}

// TestDeleteWebhook_DeleteWebhook tests the deletion of a webhook
func TestDeleteWebhook_DeleteWebhook(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/webhooks/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
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
		t.Fatalf("client creation error: %v", err)
	}

	// Delete webhook
	err = client.DeleteWebhook("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("webhook deletion error: %v", err)
	}
}
