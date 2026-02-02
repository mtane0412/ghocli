/**
 * offers.go
 * Offers API
 *
 * Provides Offers functionality for the Ghost Admin API.
 * Confirmation mechanism is applied to Create/Update operations.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Offer represents a Ghost offer
type Offer struct {
	ID                 string    `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Code               string    `json:"code,omitempty"`
	DisplayTitle       string    `json:"display_title,omitempty"`
	DisplayDescription string    `json:"display_description,omitempty"`
	Type               string    `json:"type,omitempty"`     // percent, fixed
	Cadence            string    `json:"cadence,omitempty"`  // month, year
	Amount             int       `json:"amount,omitempty"`
	Duration           string    `json:"duration,omitempty"` // once, forever, repeating
	DurationInMonths   int       `json:"duration_in_months,omitempty"`
	Currency           string    `json:"currency,omitempty"`
	Status             string    `json:"status,omitempty"`   // active, archived
	RedemptionCount    int       `json:"redemption_count,omitempty"`
	Tier               OfferTier `json:"tier,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// OfferTier represents tier information related to an offer
type OfferTier struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// OfferListOptions represents options for retrieving offer list
type OfferListOptions struct {
	Limit  int    // Number of items to retrieve (default: 15)
	Page   int    // Page number (default: 1)
	Filter string // Filter condition
}

// OfferListResponse represents an offer list response
type OfferListResponse struct {
	Offers []Offer `json:"offers"`
	Meta   struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// OfferResponse represents a single offer response
type OfferResponse struct {
	Offers []Offer `json:"offers"`
}

// ListOffers retrieves a list of offers
func (c *Client) ListOffers(opts OfferListOptions) (*OfferListResponse, error) {
	path := "/ghost/api/admin/offers/"

	// Build query parameters
	params := []string{}
	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Page > 0 {
		params = append(params, fmt.Sprintf("page=%d", opts.Page))
	}
	if opts.Filter != "" {
		params = append(params, fmt.Sprintf("filter=%s", opts.Filter))
	}

	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp OfferListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetOffer retrieves an offer by ID
// Note: Offers API does not support retrieval by slug (ID only)
func (c *Client) GetOffer(id string) (*Offer, error) {
	path := fmt.Sprintf("/ghost/api/admin/offers/%s/", id)

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("offer not found: %s", id)
	}

	return &resp.Offers[0], nil
}

// CreateOffer creates a new offer
func (c *Client) CreateOffer(offer *Offer) (*Offer, error) {
	path := "/ghost/api/admin/offers/"

	// Build request body
	reqBody := map[string]interface{}{
		"offers": []interface{}{offer},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("POST", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("failed to create offer")
	}

	return &resp.Offers[0], nil
}

// UpdateOffer updates an existing offer
func (c *Client) UpdateOffer(id string, offer *Offer) (*Offer, error) {
	path := fmt.Sprintf("/ghost/api/admin/offers/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"offers": []interface{}{offer},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("failed to update offer")
	}

	return &resp.Offers[0], nil
}
