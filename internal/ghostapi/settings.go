/**
 * settings.go
 * Settings API
 *
 * Ghost Admin APIのSettings機能を提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Setting はGhostの設定項目を表します
type Setting struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SettingsResponse は設定一覧のレスポンスです
type SettingsResponse struct {
	Settings []Setting `json:"settings"`
}

// SettingUpdate は設定更新のリクエストです
type SettingUpdate struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SettingsUpdateRequest は設定更新のリクエストボディです
type SettingsUpdateRequest struct {
	Settings []SettingUpdate `json:"settings"`
}

// GetSettings は全設定を取得します
func (c *Client) GetSettings() (*SettingsResponse, error) {
	path := "/ghost/api/admin/settings/"

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp SettingsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// UpdateSettings は設定を更新します
func (c *Client) UpdateSettings(updates []SettingUpdate) (*SettingsResponse, error) {
	path := "/ghost/api/admin/settings/"

	// リクエストボディを作成
	reqBody := SettingsUpdateRequest{
		Settings: updates,
	}

	// JSONにエンコード
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディのエンコードに失敗しました: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp SettingsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}
