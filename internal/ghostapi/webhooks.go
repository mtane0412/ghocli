/**
 * webhooks.go
 * Webhooks API
 *
 * Ghost Admin APIのWebhooks機能を提供します。
 * 注意: Ghost APIはWebhookのList/Getをサポートしていません。Create/Update/Deleteのみ提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Webhook はGhostのWebhookを表します
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

// WebhookResponse はWebhookのレスポンスです
type WebhookResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// CreateWebhook は新しいWebhookを作成します
func (c *Client) CreateWebhook(webhook *Webhook) (*Webhook, error) {
	path := "/ghost/api/admin/webhooks/"

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"webhooks": []Webhook{*webhook},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディの生成に失敗しました: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("POST", path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp WebhookResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Webhooks) == 0 {
		return nil, fmt.Errorf("Webhookの作成に失敗しました")
	}

	return &resp.Webhooks[0], nil
}

// UpdateWebhook は既存のWebhookを更新します
func (c *Client) UpdateWebhook(id string, webhook *Webhook) (*Webhook, error) {
	path := fmt.Sprintf("/ghost/api/admin/webhooks/%s/", id)

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"webhooks": []Webhook{*webhook},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディの生成に失敗しました: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp WebhookResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Webhooks) == 0 {
		return nil, fmt.Errorf("Webhookの更新に失敗しました")
	}

	return &resp.Webhooks[0], nil
}

// DeleteWebhook はWebhookを削除します
func (c *Client) DeleteWebhook(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/webhooks/%s/", id)

	// リクエストを実行
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
