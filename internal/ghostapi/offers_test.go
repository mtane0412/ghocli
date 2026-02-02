/**
 * offers_test.go
 * Test code for Offers API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListOffers_GetOfferList retrieves a list of offers
func TestListOffers_GetOfferList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		if r.URL.Path != "/ghost/api/admin/offers/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/offers/")
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
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "Spring Campaign",
					"code":               "SPRING2024",
					"display_title":       "Spring Special Discount",
					"display_description": "50% off for a limited time",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "Premium Membership",
					},
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":                  "64fac5417c4c6b0001234568",
					"name":               "New Member Benefit",
					"code":               "WELCOME100",
					"display_title":       "100 yen off for new registration",
					"display_description": "First month only",
					"type":               "fixed",
					"cadence":            "month",
					"amount":             100,
					"duration":           "once",
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   25,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "Premium Membership",
					},
					"created_at": "2024-01-16T10:00:00.000Z",
					"updated_at": "2024-01-16T10:00:00.000Z",
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
		t.Fatalf("Client creation error: %v", err)
	}

	// Get offer list
	resp, err := client.ListOffers(OfferListOptions{})
	if err != nil {
		t.Fatalf("Offer list retrieval error: %v", err)
	}

	// Validate response
	if len(resp.Offers) != 2 {
		t.Errorf("Number of offers = %d; want 2", len(resp.Offers))
	}

	// Validate first offer
	firstOffer := resp.Offers[0]
	if firstOffer.Name != "Spring Campaign" {
		t.Errorf("Offer name = %q; want %q", firstOffer.Name, "Spring Campaign")
	}
	if firstOffer.Code != "SPRING2024" {
		t.Errorf("Code = %q; want %q", firstOffer.Code, "SPRING2024")
	}
	if firstOffer.Type != "percent" {
		t.Errorf("Type = %q; want %q", firstOffer.Type, "percent")
	}
}

// TestListOffers_FilterParameter tests filter parameter
func TestListOffers_FilterParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate query parameter
		if !r.URL.Query().Has("filter") {
			t.Error("filter parameter is not set")
		}
		if r.URL.Query().Get("filter") != "status:active" {
			t.Errorf("filter parameter = %q; want %q", r.URL.Query().Get("filter"), "status:active")
		}

		// Return response
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "Spring Campaign",
					"code":               "SPRING2024",
					"display_title":       "Spring Special Discount",
					"display_description": "50% off for a limited time",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"created_at":         "2024-01-15T10:00:00.000Z",
					"updated_at":         "2024-01-15T10:00:00.000Z",
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
		t.Fatalf("Client creation error: %v", err)
	}

	// Get offer list (filtered by status:active)
	resp, err := client.ListOffers(OfferListOptions{
		Filter: "status:active",
	})
	if err != nil {
		t.Fatalf("Offer list retrieval error: %v", err)
	}

	// Validate response
	if len(resp.Offers) != 1 {
		t.Errorf("Number of offers = %d; want 1", len(resp.Offers))
	}
}

// TestGetOffer_GetOfferByID retrieves an offer by ID
func TestGetOffer_GetOfferByID(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		expectedPath := "/ghost/api/admin/offers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "Spring Campaign",
					"code":               "SPRING2024",
					"display_title":       "Spring Special Discount",
					"display_description": "50% off for a limited time",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "Premium Membership",
					},
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
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
		t.Fatalf("Client creation error: %v", err)
	}

	// Get offer
	offer, err := client.GetOffer("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Offer retrieval error: %v", err)
	}

	// Validate response
	if offer.Name != "Spring Campaign" {
		t.Errorf("Offer name = %q; want %q", offer.Name, "Spring Campaign")
	}
	if offer.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("Offer ID = %q; want %q", offer.ID, "64fac5417c4c6b0001234567")
	}
}

// TestCreateOffer_CreateOffer creates a new offer
func TestCreateOffer_CreateOffer(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		if r.URL.Path != "/ghost/api/admin/offers/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/offers/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Validate request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		offers, ok := reqBody["offers"].([]interface{})
		if !ok || len(offers) == 0 {
			t.Error("offers field is incorrect")
		}

		// Return response
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234569",
					"name":               "New Offer",
					"code":               "NEWCODE",
					"display_title":       "Test Offer",
					"display_description": "Test description",
					"type":               "percent",
					"cadence":            "month",
					"amount":             30,
					"duration":           "once",
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   0,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "Premium Membership",
					},
					"created_at": "2024-01-20T10:00:00.000Z",
					"updated_at": "2024-01-20T10:00:00.000Z",
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
		t.Fatalf("Client creation error: %v", err)
	}

	// Create offer
	newOffer := &Offer{
		Name:               "New Offer",
		Code:               "NEWCODE",
		DisplayTitle:       "Test Offer",
		DisplayDescription: "Test description",
		Type:               "percent",
		Cadence:            "month",
		Amount:             30,
		Duration:           "once",
		Currency:           "JPY",
		Tier: OfferTier{
			ID: "64fac5417c4c6b0001234999",
		},
	}

	createdOffer, err := client.CreateOffer(newOffer)
	if err != nil {
		t.Fatalf("Offer creation error: %v", err)
	}

	// Validate response
	if createdOffer.Name != "New Offer" {
		t.Errorf("Offer name = %q; want %q", createdOffer.Name, "New Offer")
	}
	if createdOffer.ID == "" {
		t.Error("Offer ID is not set")
	}
}

// TestUpdateOffer_UpdateOffer updates an existing offer
func TestUpdateOffer_UpdateOffer(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		expectedPath := "/ghost/api/admin/offers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Validate request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		offers, ok := reqBody["offers"].([]interface{})
		if !ok || len(offers) == 0 {
			t.Error("offers field is incorrect")
		}

		// Return response
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "Updated Offer",
					"code":               "SPRING2024",
					"display_title":       "Updated Title",
					"display_description": "Updated description",
					"type":               "percent",
					"cadence":            "month",
					"amount":             60,
					"duration":           "repeating",
					"duration_in_months": 6,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   15,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "Premium Membership",
					},
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-20T10:00:00.000Z",
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
		t.Fatalf("Client creation error: %v", err)
	}

	// Update offer
	updateOffer := &Offer{
		Name:               "Updated Offer",
		DisplayTitle:       "Updated Title",
		DisplayDescription: "Updated description",
		Amount:             60,
		DurationInMonths:   6,
	}

	updatedOffer, err := client.UpdateOffer("64fac5417c4c6b0001234567", updateOffer)
	if err != nil {
		t.Fatalf("Offer update error: %v", err)
	}

	// Validate response
	if updatedOffer.Name != "Updated Offer" {
		t.Errorf("Offer name = %q; want %q", updatedOffer.Name, "Updated Offer")
	}
	if updatedOffer.DisplayTitle != "Updated Title" {
		t.Errorf("Display title = %q; want %q", updatedOffer.DisplayTitle, "Updated Title")
	}
}
