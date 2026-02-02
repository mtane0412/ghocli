/**
 * settings.go
 * Settings API
 *
 * Provides Settings functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Setting represents a Ghost settings item
type Setting struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SettingsResponse represents a settings list response
type SettingsResponse struct {
	Settings []Setting `json:"settings"`
}

// SettingUpdate represents a settings update request
type SettingUpdate struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SettingsUpdateRequest represents a settings update request body
type SettingsUpdateRequest struct {
	Settings []SettingUpdate `json:"settings"`
}

// GetSettings retrieves all settings
func (c *Client) GetSettings() (*SettingsResponse, error) {
	path := "/ghost/api/admin/settings/"

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp SettingsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// UpdateSettings updates settings
func (c *Client) UpdateSettings(updates []SettingUpdate) (*SettingsResponse, error) {
	path := "/ghost/api/admin/settings/"

	// Create request body
	reqBody := SettingsUpdateRequest{
		Settings: updates,
	}

	// Encode to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp SettingsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}
