/**
 * webhooks.go
 * Webhooks API
 *
 * Provides Webhooks functionality for the Ghost Admin API.
 * Note: Ghost API does not support List/Get operations for Webhooks. Only Create/Update/Delete are provided.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Webhook represents a Ghost webhook
type Webhook struct {
	ID              string     `json:"id,omitempty"`
	Event           string     `json:"event"`
	TargetURL       string     `json:"target_url"`
	Name            string     `json:"name,omitempty"`
	Secret          string     `json:"secret,omitempty"`
	APIVersion      string     `json:"api_version,omitempty"`
	IntegrationID   string     `json:"integration_id,omitempty"`
	Status          string     `json:"status,omitempty"`
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
}

// WebhookResponse represents a webhook response
type WebhookResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(webhook *Webhook) (*Webhook, error) {
	path := "/ghost/api/admin/webhooks/"

	// Build request body
	reqBody := map[string]interface{}{
		"webhooks": []Webhook{*webhook},
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
	var resp WebhookResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Webhooks) == 0 {
		return nil, fmt.Errorf("failed to create webhook")
	}

	return &resp.Webhooks[0], nil
}

// UpdateWebhook updates an existing webhook
func (c *Client) UpdateWebhook(id string, webhook *Webhook) (*Webhook, error) {
	path := fmt.Sprintf("/ghost/api/admin/webhooks/%s/", id)

	// Build request body
	reqBody := map[string]interface{}{
		"webhooks": []Webhook{*webhook},
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
	var resp WebhookResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Webhooks) == 0 {
		return nil, fmt.Errorf("failed to update webhook")
	}

	return &resp.Webhooks[0], nil
}

// DeleteWebhook deletes a webhook
func (c *Client) DeleteWebhook(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/webhooks/%s/", id)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
