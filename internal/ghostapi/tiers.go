/**
 * tiers.go
 * Tiers API
 *
 * Provides Tiers functionality for the Ghost Admin API.
 * Create/Update operations are subject to confirmation mechanisms.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tier represents a Ghost tier
type Tier struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	Active         bool      `json:"active,omitempty"`
	Type           string    `json:"type,omitempty"`          // free, paid
	Visibility     string    `json:"visibility,omitempty"`    // public, none
	WelcomePageURL string    `json:"welcome_page_url,omitempty"`
	MonthlyPrice   int       `json:"monthly_price,omitempty"` // smallest currency unit
	YearlyPrice    int       `json:"yearly_price,omitempty"`
	Currency       string    `json:"currency,omitempty"`
	Benefits       []string  `json:"benefits,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// TierListOptions represents options for fetching tier list
type TierListOptions struct {
	Limit   int    // Number of items to fetch (default: 15)
	Page    int    // Page number (default: 1)
	Include string // Additional information to include (monthly_price, yearly_price, benefits, etc.)
	Filter  string // Filter condition
}

// TierListResponse represents a tier list response
type TierListResponse struct {
	Tiers []Tier `json:"tiers"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// TierResponse represents a single tier response
type TierResponse struct {
	Tiers []Tier `json:"tiers"`
}

// ListTiers retrieves a list of tiers
func (c *Client) ListTiers(opts TierListOptions) (*TierListResponse, error) {
	path := "/ghost/api/admin/tiers/"

	// Build query parameters
	params := []string{}
	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Page > 0 {
		params = append(params, fmt.Sprintf("page=%d", opts.Page))
	}
	if opts.Include != "" {
		params = append(params, fmt.Sprintf("include=%s", opts.Include))
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
	var resp TierListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetTier retrieves a tier by ID or slug
// If idOrSlug starts with "slug:", it will be treated as a slug
func (c *Client) GetTier(idOrSlug string) (*Tier, error) {
	var path string

	// Determine if it's a slug or ID
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/tiers/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/tiers/%s/", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TierResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tiers) == 0 {
		return nil, fmt.Errorf("tier not found: %s", idOrSlug)
	}

	return &resp.Tiers[0], nil
}

// CreateTier creates a new tier
func (c *Client) CreateTier(tier *Tier) (*Tier, error) {
	path := "/ghost/api/admin/tiers/"

	// Build request body
	reqBody := map[string]interface{}{
		"tiers": []interface{}{tier},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("POST", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TierResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tiers) == 0 {
		return nil, fmt.Errorf("failed to create tier")
	}

	return &resp.Tiers[0], nil
}

// UpdateTier updates an existing tier
func (c *Client) UpdateTier(id string, tier *Tier) (*Tier, error) {
	path := fmt.Sprintf("/ghost/api/admin/tiers/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"tiers": []interface{}{tier},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TierResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tiers) == 0 {
		return nil, fmt.Errorf("failed to update tier")
	}

	return &resp.Tiers[0], nil
}
