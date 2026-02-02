/**
 * client_test.go
 * Test code for Ghost Admin API client
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewClient_CreateClient creates a client
func TestNewClient_CreateClient(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	siteURL := "https://test.ghost.io"

	client, err := NewClient(siteURL, keyID, secret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if client.baseURL != siteURL {
		t.Errorf("baseURL = %q; want %q", client.baseURL, siteURL)
	}
}

// TestNewClient_ErrorWithInvalidURL tests error with invalid URL
func TestNewClient_ErrorWithInvalidURL(t *testing.T) {
	_, err := NewClient("", "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		t.Error("No error returned with empty URL")
	}
}

// TestGetSite_RetrieveSiteInformation retrieves site information
func TestGetSite_RetrieveSiteInformation(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/site/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/site/")
		}

		// Verify Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}
		if len(auth) < 6 || auth[:6] != "Ghost " {
			t.Errorf("Authorization header is invalid: %s", auth)
		}

		// Return response
		response := map[string]interface{}{
			"site": map[string]interface{}{
				"title":       "Test Blog",
				"description": "A test blog",
				"url":         "https://test.ghost.io",
				"version":     "5.0",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Retrieve site information
	site, err := client.GetSite()
	if err != nil {
		t.Fatalf("Failed to retrieve site information: %v", err)
	}

	// Verify response
	if site.Title != "Test Blog" {
		t.Errorf("Title = %q; want %q", site.Title, "Test Blog")
	}
	if site.Description != "A test blog" {
		t.Errorf("Description = %q; want %q", site.Description, "A test blog")
	}
	if site.URL != "https://test.ghost.io" {
		t.Errorf("URL = %q; want %q", site.URL, "https://test.ghost.io")
	}
	if site.Version != "5.0" {
		t.Errorf("Version = %q; want %q", site.Version, "5.0")
	}
}

// TestGetSite_APIError tests API error
func TestGetSite_APIError(t *testing.T) {
	// Create HTTP server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"message": "Unauthorized",
					"type":    "UnauthorizedError",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "invalid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Retrieve site information (expect error)
	_, err = client.GetSite()
	if err == nil {
		t.Error("expected error but got nil")
	}
}

// TestDoRequestWithOptions_QueryParamsWithJapaneseAreEncoded tests Japanese query parameters are encoded
func TestDoRequestWithOptions_QueryParamsWithJapaneseAreEncoded(t *testing.T) {
	// Create HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		title := query.Get("title")
		if title != "テスト投稿" {
			t.Errorf("title = %q; want %q", title, "テスト投稿")
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Execute request with query parameters
	opts := &RequestOptions{
		QueryParams: map[string]string{
			"title": "テスト投稿",
		},
	}
	_, err = client.doRequestWithOptions("GET", "/test", nil, opts)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
}

// TestDoRequestWithOptions_QueryParamsWithAmpersandAreEncoded tests ampersand query parameters are encoded
func TestDoRequestWithOptions_QueryParamsWithAmpersandAreEncoded(t *testing.T) {
	// Create HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		value := query.Get("value")
		if value != "foo&bar" {
			t.Errorf("value = %q; want %q", value, "foo&bar")
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Execute request with query parameters
	opts := &RequestOptions{
		QueryParams: map[string]string{
			"value": "foo&bar",
		},
	}
	_, err = client.doRequestWithOptions("GET", "/test", nil, opts)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
}

// TestDoRequestWithOptions_QueryParamsWithEqualsAreEncoded tests equals query parameters are encoded
func TestDoRequestWithOptions_QueryParamsWithEqualsAreEncoded(t *testing.T) {
	// Create HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		value := query.Get("value")
		if value != "foo=bar" {
			t.Errorf("value = %q; want %q", value, "foo=bar")
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Execute request with query parameters
	opts := &RequestOptions{
		QueryParams: map[string]string{
			"value": "foo=bar",
		},
	}
	_, err = client.doRequestWithOptions("GET", "/test", nil, opts)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
}

// TestDoRequestWithOptions_QueryParamsWithSpaceAreEncoded tests space query parameters are encoded
func TestDoRequestWithOptions_QueryParamsWithSpaceAreEncoded(t *testing.T) {
	// Create HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		value := query.Get("value")
		if value != "foo bar" {
			t.Errorf("value = %q; want %q", value, "foo bar")
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Execute request with query parameters
	opts := &RequestOptions{
		QueryParams: map[string]string{
			"value": "foo bar",
		},
	}
	_, err = client.doRequestWithOptions("GET", "/test", nil, opts)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
}

// TestDoRequestWithOptions_MultipleQueryParamsAreEncoded tests multiple query parameters are encoded
func TestDoRequestWithOptions_MultipleQueryParamsAreEncoded(t *testing.T) {
	// Create HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		name := query.Get("name")
		value := query.Get("value")
		if name != "John Doe" {
			t.Errorf("name = %q; want %q", name, "John Doe")
		}
		if value != "foo&bar=baz" {
			t.Errorf("value = %q; want %q", value, "foo&bar=baz")
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Execute request with multiple query parameters
	opts := &RequestOptions{
		QueryParams: map[string]string{
			"name":  "John Doe",
			"value": "foo&bar=baz",
		},
	}
	_, err = client.doRequestWithOptions("GET", "/test", nil, opts)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
}
