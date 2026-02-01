/**
 * themes.go
 * Themes API
 *
 * Ghost Admin APIのThemes機能を提供します。
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// Theme はGhostのテーマを表します
type Theme struct {
	Name      string          `json:"name"`
	Package   *ThemePackage   `json:"package,omitempty"`
	Active    bool            `json:"active"`
	Templates []ThemeTemplate `json:"templates,omitempty"`
}

// ThemePackage はテーマのパッケージ情報を表します
type ThemePackage struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

// ThemeTemplate はテーマのテンプレートファイル情報を表します
type ThemeTemplate struct {
	Filename string `json:"filename"`
}

// ThemeListResponse はテーマ一覧のレスポンスです
type ThemeListResponse struct {
	Themes []Theme `json:"themes"`
}

// ThemeResponse はテーマ単体のレスポンスです
type ThemeResponse struct {
	Themes []Theme `json:"themes"`
}

// ListThemes はテーマ一覧を取得します
func (c *Client) ListThemes() (*ThemeListResponse, error) {
	path := "/ghost/api/admin/themes/"

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp ThemeListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// UploadTheme はテーマをアップロードします
func (c *Client) UploadTheme(file io.Reader, filename string) (*Theme, error) {
	path := "/ghost/api/admin/themes/upload/"

	// マルチパートリクエストを実行
	respBody, err := c.doMultipartRequest(path, file, filename, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp ThemeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Themes) == 0 {
		return nil, fmt.Errorf("テーマのアップロードに失敗しました")
	}

	return &resp.Themes[0], nil
}

// ActivateTheme はテーマを有効化します
func (c *Client) ActivateTheme(name string) (*Theme, error) {
	path := fmt.Sprintf("/ghost/api/admin/themes/%s/activate/", name)

	// リクエストを実行
	respBody, err := c.doRequest("PUT", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp ThemeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Themes) == 0 {
		return nil, fmt.Errorf("テーマの有効化に失敗しました")
	}

	return &resp.Themes[0], nil
}

// DeleteTheme はテーマを削除します
func (c *Client) DeleteTheme(name string) error {
	path := fmt.Sprintf("/ghost/api/admin/themes/%s/", name)

	// リクエストを実行
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
