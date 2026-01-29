/**
 * tags.go
 * Tags API
 *
 * Ghost Admin APIのTags機能を提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tag はGhostのタグを表します
type Tag struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	Visibility  string    `json:"visibility,omitempty"` // public, internal
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// TagListOptions はタグ一覧取得のオプションです
type TagListOptions struct {
	Limit   int    // 取得件数（デフォルト: 15）
	Page    int    // ページ番号（デフォルト: 1）
	Include string // 含める追加情報（count.posts など）
	Filter  string // フィルター条件
}

// TagListResponse はタグ一覧のレスポンスです
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

// TagResponse はタグ単体のレスポンスです
type TagResponse struct {
	Tags []Tag `json:"tags"`
}

// ListTags はタグ一覧を取得します
func (c *Client) ListTags(opts TagListOptions) (*TagListResponse, error) {
	path := "/ghost/api/admin/tags/"

	// クエリパラメータを構築
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

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp TagListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetTag は指定されたIDまたはスラッグのタグを取得します
// idOrSlugが "slug:" で始まる場合はスラッグとして扱います
func (c *Client) GetTag(idOrSlug string) (*Tag, error) {
	var path string

	// スラッグかIDかを判定
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/tags/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/tags/%s/", idOrSlug)
	}

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("タグが見つかりません: %s", idOrSlug)
	}

	return &resp.Tags[0], nil
}

// CreateTag は新しいタグを作成します
func (c *Client) CreateTag(tag *Tag) (*Tag, error) {
	path := "/ghost/api/admin/tags/"

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"tags": []Tag{*tag},
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
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("タグの作成に失敗しました")
	}

	return &resp.Tags[0], nil
}

// UpdateTag は既存のタグを更新します
func (c *Client) UpdateTag(id string, tag *Tag) (*Tag, error) {
	path := fmt.Sprintf("/ghost/api/admin/tags/%s/", id)

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"tags": []Tag{*tag},
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
	var resp TagResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("タグの更新に失敗しました")
	}

	return &resp.Tags[0], nil
}

// DeleteTag はタグを削除します
func (c *Client) DeleteTag(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/tags/%s/", id)

	// リクエストを実行
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
