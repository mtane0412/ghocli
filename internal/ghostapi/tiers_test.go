/**
 * tiers_test.go
 * Test code for Tiers API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListTiers_GetTierList tests fetching a list of tiers
func TestListTiers_GetTierList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request
		if r.URL.Path != "/ghost/api/admin/tiers/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/tiers/")
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
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "Free Member",
					"description":      "Read articles for free",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
				},
				{
					"id":               "64fac5417c4c6b0001234568",
					"name":             "Premium Member",
					"description":      "Access to all articles",
					"slug":             "premium",
					"active":           true,
					"type":             "paid",
					"visibility":       "public",
					"monthly_price":    500,
					"yearly_price":     5000,
					"currency":         "JPY",
					"welcome_page_url": "/welcome",
					"created_at":       "2024-01-16T10:00:00.000Z",
					"updated_at":       "2024-01-16T10:00:00.000Z",
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
		t.Fatalf("client creation error: %v", err)
	}

	// Get tier list
	resp, err := client.ListTiers(TierListOptions{})
	if err != nil {
		t.Fatalf("tier list retrieval error: %v", err)
	}

	// Validate response
	if len(resp.Tiers) != 2 {
		t.Errorf("number of tiers = %d; want 2", len(resp.Tiers))
	}

	// Validate first tier
	firstTier := resp.Tiers[0]
	if firstTier.Name != "Free Member" {
		t.Errorf("tier name = %q; want %q", firstTier.Name, "Free Member")
	}
	if firstTier.Slug != "free" {
		t.Errorf("slug = %q; want %q", firstTier.Slug, "free")
	}
	if firstTier.Type != "free" {
		t.Errorf("type = %q; want %q", firstTier.Type, "free")
	}
}

// TestListTiers_IncludeParameter tests the include parameter
func TestListTiers_IncludeParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate query parameters
		if !r.URL.Query().Has("include") {
			t.Error("include parameter is not set")
		}
		if r.URL.Query().Get("include") != "monthly_price,yearly_price" {
			t.Errorf("include parameter = %q; want %q", r.URL.Query().Get("include"), "monthly_price,yearly_price")
		}

		// Return response
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":            "64fac5417c4c6b0001234568",
					"name":          "Premium Member",
					"slug":          "premium",
					"type":          "paid",
					"active":        true,
					"visibility":    "public",
					"monthly_price": 500,
					"yearly_price":  5000,
					"currency":      "JPY",
					"created_at":    "2024-01-16T10:00:00.000Z",
					"updated_at":    "2024-01-16T10:00:00.000Z",
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
		t.Fatalf("client creation error: %v", err)
	}

	// Get tier list (including monthly_price, yearly_price)
	resp, err := client.ListTiers(TierListOptions{
		Include: "monthly_price,yearly_price",
	})
	if err != nil {
		t.Fatalf("tier list retrieval error: %v", err)
	}

	// Validate response
	if len(resp.Tiers) != 1 {
		t.Errorf("number of tiers = %d; want 1", len(resp.Tiers))
	}
}

// TestGetTier_GetTierByID tests fetching a tier by ID
func TestGetTier_GetTierByID(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request
		expectedPath := "/ghost/api/admin/tiers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "Free Member",
					"description":      "Read articles for free",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
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

	// Get tier
	tier, err := client.GetTier("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("tier retrieval error: %v", err)
	}

	// Validate response
	if tier.Name != "Free Member" {
		t.Errorf("tier name = %q; want %q", tier.Name, "Free Member")
	}
	if tier.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("tier ID = %q; want %q", tier.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetTier_GetTierBySlug tests fetching a tier by slug
func TestGetTier_GetTierBySlug(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request
		expectedPath := "/ghost/api/admin/tiers/slug/free/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "Free Member",
					"description":      "Read articles for free",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
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

	// Get tier
	tier, err := client.GetTier("slug:free")
	if err != nil {
		t.Fatalf("tier retrieval error: %v", err)
	}

	// Validate response
	if tier.Slug != "free" {
		t.Errorf("slug = %q; want %q", tier.Slug, "free")
	}
}

// TestCreateTier_CreateTier tests creating a tier
func TestCreateTier_CreateTier(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request
		if r.URL.Path != "/ghost/api/admin/tiers/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/tiers/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Validate request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		tiers, ok := reqBody["tiers"].([]interface{})
		if !ok || len(tiers) == 0 {
			t.Error("tiers field is incorrect")
		}

		// Return response
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234569",
					"name":             "New Plan",
					"description":      "Test Plan",
					"slug":             "new-plan",
					"active":           true,
					"type":             "paid",
					"visibility":       "public",
					"monthly_price":    1000,
					"yearly_price":     10000,
					"currency":         "JPY",
					"welcome_page_url": "/welcome",
					"created_at":       "2024-01-20T10:00:00.000Z",
					"updated_at":       "2024-01-20T10:00:00.000Z",
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

	// Create tier
	newTier := &Tier{
		Name:           "New Plan",
		Description:    "Test Plan",
		Type:           "paid",
		Visibility:     "public",
		MonthlyPrice:   1000,
		YearlyPrice:    10000,
		Currency:       "JPY",
		WelcomePageURL: "/welcome",
	}

	createdTier, err := client.CreateTier(newTier)
	if err != nil {
		t.Fatalf("tier creation error: %v", err)
	}

	// Validate response
	if createdTier.Name != "New Plan" {
		t.Errorf("tier name = %q; want %q", createdTier.Name, "New Plan")
	}
	if createdTier.ID == "" {
		t.Error("tier ID is not set")
	}
}

// TestUpdateTier_UpdateTier tests updating a tier
func TestUpdateTier_UpdateTier(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request
		expectedPath := "/ghost/api/admin/tiers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Validate request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		tiers, ok := reqBody["tiers"].([]interface{})
		if !ok || len(tiers) == 0 {
			t.Error("tiers field is incorrect")
		}

		// Return response
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "Updated Plan",
					"description":      "Updated description",
					"slug":             "premium",
					"active":           true,
					"type":             "paid",
					"visibility":       "public",
					"monthly_price":    1500,
					"yearly_price":     15000,
					"currency":         "JPY",
					"welcome_page_url": "/updated-welcome",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-20T10:00:00.000Z",
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

	// Update tier
	updateTier := &Tier{
		Name:           "Updated Plan",
		Description:    "Updated description",
		MonthlyPrice:   1500,
		YearlyPrice:    15000,
		WelcomePageURL: "/updated-welcome",
	}

	updatedTier, err := client.UpdateTier("64fac5417c4c6b0001234567", updateTier)
	if err != nil {
		t.Fatalf("tier update error: %v", err)
	}

	// Validate response
	if updatedTier.Name != "Updated Plan" {
		t.Errorf("tier name = %q; want %q", updatedTier.Name, "Updated Plan")
	}
	if updatedTier.Description != "Updated description" {
		t.Errorf("description = %q; want %q", updatedTier.Description, "Updated description")
	}
}
