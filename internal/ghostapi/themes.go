/**
 * themes.go
 * Themes API
 *
 * Provides Themes functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// Theme represents a Ghost theme
type Theme struct {
	Name      string          `json:"name"`
	Package   *ThemePackage   `json:"package,omitempty"`
	Active    bool            `json:"active"`
	Templates []ThemeTemplate `json:"templates,omitempty"`
}

// ThemePackage represents theme package information
type ThemePackage struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

// ThemeTemplate represents theme template file information
type ThemeTemplate struct {
	Filename string `json:"filename"`
}

// ThemeListResponse represents a theme list response
type ThemeListResponse struct {
	Themes []Theme `json:"themes"`
}

// ThemeResponse represents a single theme response
type ThemeResponse struct {
	Themes []Theme `json:"themes"`
}

// ListThemes retrieves a list of themes
func (c *Client) ListThemes() (*ThemeListResponse, error) {
	path := "/ghost/api/admin/themes/"

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ThemeListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// UploadTheme uploads a theme
func (c *Client) UploadTheme(file io.Reader, filename string) (*Theme, error) {
	path := "/ghost/api/admin/themes/upload/"

	// Execute multipart request
	respBody, err := c.doMultipartRequest(path, file, filename, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ThemeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Themes) == 0 {
		return nil, fmt.Errorf("failed to upload theme")
	}

	return &resp.Themes[0], nil
}

// ActivateTheme activates a theme
func (c *Client) ActivateTheme(name string) (*Theme, error) {
	path := fmt.Sprintf("/ghost/api/admin/themes/%s/activate/", name)

	// Execute request
	respBody, err := c.doRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ThemeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Themes) == 0 {
		return nil, fmt.Errorf("failed to activate theme")
	}

	return &resp.Themes[0], nil
}

// DeleteTheme deletes a theme
func (c *Client) DeleteTheme(name string) error {
	path := fmt.Sprintf("/ghost/api/admin/themes/%s/", name)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
