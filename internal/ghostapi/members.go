/**
 * members.go
 * Members API
 *
 * Provides Members functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Member represents a Ghost member (subscriber)
type Member struct {
	ID        string    `json:"id,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	Email     string    `json:"email"`              // Required field
	Name      string    `json:"name,omitempty"`
	Note      string    `json:"note,omitempty"`
	Status    string    `json:"status,omitempty"`   // free, paid, comped
	Labels    []Label   `json:"labels,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Label represents a label assigned to a member
type Label struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// MemberListOptions represents options for fetching member list
type MemberListOptions struct {
	Limit  int    // Number of items to fetch (default: 15)
	Page   int    // Page number (default: 1)
	Filter string // Filter condition
	Order  string // Sort order
}

// MemberListResponse represents a member list response
type MemberListResponse struct {
	Members []Member `json:"members"`
	Meta    struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// MemberResponse represents a single member response
type MemberResponse struct {
	Members []Member `json:"members"`
}

// ListMembers retrieves a list of members
func (c *Client) ListMembers(opts MemberListOptions) (*MemberListResponse, error) {
	path := "/ghost/api/admin/members/"

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
	if opts.Order != "" {
		params = append(params, fmt.Sprintf("order=%s", opts.Order))
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
	var resp MemberListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetMember retrieves a member by ID
func (c *Client) GetMember(id string) (*Member, error) {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("member not found: %s", id)
	}

	return &resp.Members[0], nil
}

// CreateMember creates a new member
func (c *Client) CreateMember(member *Member) (*Member, error) {
	path := "/ghost/api/admin/members/"

	// Build request body
	reqBody := map[string]interface{}{
		"members": []Member{*member},
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
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("failed to create member")
	}

	return &resp.Members[0], nil
}

// UpdateMember updates an existing member
func (c *Client) UpdateMember(id string, member *Member) (*Member, error) {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"members": []Member{*member},
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
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("failed to update member")
	}

	return &resp.Members[0], nil
}

// DeleteMember deletes a member
func (c *Client) DeleteMember(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
