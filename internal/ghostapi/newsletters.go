/**
 * newsletters.go
 * Newsletters API
 *
 * Provides Newsletters functionality for the Ghost Admin API.
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

// Newsletter represents a Ghost newsletter
type Newsletter struct {
	ID                string    `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	Slug              string    `json:"slug,omitempty"`
	Status            string    `json:"status,omitempty"`            // active, archived
	Visibility        string    `json:"visibility,omitempty"`        // members, paid
	SubscribeOnSignup bool      `json:"subscribe_on_signup,omitempty"`
	SenderName        string    `json:"sender_name,omitempty"`
	SenderEmail       string    `json:"sender_email,omitempty"`
	SenderReplyTo     string    `json:"sender_reply_to,omitempty"`
	SortOrder         int       `json:"sort_order,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

// NewsletterListOptions represents options for retrieving newsletter list
type NewsletterListOptions struct {
	Limit  int    // Number of items to retrieve (default: 15)
	Page   int    // Page number (default: 1)
	Filter string // Filter condition
}

// NewsletterListResponse represents a newsletter list response
type NewsletterListResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
	Meta        struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// NewsletterResponse represents a single newsletter response
type NewsletterResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
}

// ListNewsletters retrieves a list of newsletters
func (c *Client) ListNewsletters(opts NewsletterListOptions) (*NewsletterListResponse, error) {
	path := "/ghost/api/admin/newsletters/"

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
	var resp NewsletterListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetNewsletter retrieves a newsletter by ID or slug
// If idOrSlug starts with "slug:", it is treated as a slug
func (c *Client) GetNewsletter(idOrSlug string) (*Newsletter, error) {
	var path string

	// Determine if it's a slug or ID
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/newsletters/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/newsletters/%s/", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp NewsletterResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Newsletters) == 0 {
		return nil, fmt.Errorf("newsletter not found: %s", idOrSlug)
	}

	return &resp.Newsletters[0], nil
}

// CreateNewsletter creates a new newsletter
func (c *Client) CreateNewsletter(newsletter *Newsletter) (*Newsletter, error) {
	path := "/ghost/api/admin/newsletters/"

	// Build request body
	reqBody := map[string]interface{}{
		"newsletters": []interface{}{newsletter},
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
	var resp NewsletterResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Newsletters) == 0 {
		return nil, fmt.Errorf("failed to create newsletter")
	}

	return &resp.Newsletters[0], nil
}

// UpdateNewsletter updates an existing newsletter
func (c *Client) UpdateNewsletter(id string, newsletter *Newsletter) (*Newsletter, error) {
	path := fmt.Sprintf("/ghost/api/admin/newsletters/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"newsletters": []interface{}{newsletter},
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
	var resp NewsletterResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Newsletters) == 0 {
		return nil, fmt.Errorf("failed to update newsletter")
	}

	return &resp.Newsletters[0], nil
}
