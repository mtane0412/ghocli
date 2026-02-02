/**
 * users.go
 * Users API
 *
 * Provides Users functionality for Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User represents a Ghost user
type User struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Slug         string    `json:"slug,omitempty"`
	Email        string    `json:"email,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Location     string    `json:"location,omitempty"`
	Website      string    `json:"website,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	CoverImage   string    `json:"cover_image,omitempty"`
	Roles        []Role    `json:"roles,omitempty"` // Read-only
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Role represents a user role
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserListOptions contains options for listing users
type UserListOptions struct {
	Limit   int    // Number of users to retrieve (default: 15)
	Page    int    // Page number (default: 1)
	Include string // Additional information to include (roles, count.posts, etc.)
	Filter  string // Filter condition
}

// UserListResponse represents the response of user list
type UserListResponse struct {
	Users []User `json:"users"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// UserResponse represents the response of a single user
type UserResponse struct {
	Users []User `json:"users"`
}

// ListUsers retrieves a list of users
func (c *Client) ListUsers(opts UserListOptions) (*UserListResponse, error) {
	path := "/ghost/api/admin/users/"

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
	var resp UserListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetUser retrieves a user by ID or slug
// If idOrSlug starts with "slug:", it is treated as a slug
func (c *Client) GetUser(idOrSlug string) (*User, error) {
	var path string

	// Determine whether it's a slug or ID
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/users/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/users/%s/", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp UserResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("user not found: %s", idOrSlug)
	}

	return &resp.Users[0], nil
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(id string, user *User) (*User, error) {
	path := fmt.Sprintf("/ghost/api/admin/users/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"users": []User{*user},
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
	var resp UserResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("failed to update user")
	}

	return &resp.Users[0], nil
}
