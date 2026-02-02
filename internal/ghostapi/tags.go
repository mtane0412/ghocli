/**
 * tags.go
 * Tags API
 *
 * Provides Tags functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tag represents a Ghost tag
type Tag struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	Visibility  string    `json:"visibility,omitempty"` // public, internal
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// TagListOptions represents options for fetching tag list
type TagListOptions struct {
	Limit   int    // Number of items to fetch (default: 15)
	Page    int    // Page number (default: 1)
	Include string // Additional information to include (count.posts, etc.)
	Filter  string // Filter condition
}

// TagListResponse represents a tag list response
type TagListResponse struct {
	Tags []Tag `json:"tags"`
	Meta struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// TagResponse represents a single tag response
type TagResponse struct {
	Tags []Tag `json:"tags"`
}

// ListTags retrieves a list of tags
func (c *Client) ListTags(opts TagListOptions) (*TagListResponse, error) {
	path := "/ghost/api/admin/tags/"

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
	var resp TagListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetTag retrieves a tag by ID or slug
// If idOrSlug starts with "slug:", it will be treated as a slug
func (c *Client) GetTag(idOrSlug string) (*Tag, error) {
	var path string

	// Determine if it's a slug or ID
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/tags/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/tags/%s/", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("tag not found: %s", idOrSlug)
	}

	return &resp.Tags[0], nil
}

// CreateTag creates a new tag
func (c *Client) CreateTag(tag *Tag) (*Tag, error) {
	path := "/ghost/api/admin/tags/"

	// Build request body
	reqBody := map[string]interface{}{
		"tags": []Tag{*tag},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("POST", path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("failed to create tag")
	}

	return &resp.Tags[0], nil
}

// UpdateTag updates an existing tag
func (c *Client) UpdateTag(id string, tag *Tag) (*Tag, error) {
	path := fmt.Sprintf("/ghost/api/admin/tags/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"tags": []Tag{*tag},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("failed to update tag")
	}

	return &resp.Tags[0], nil
}

// DeleteTag deletes a tag
func (c *Client) DeleteTag(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/tags/%s/", id)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
