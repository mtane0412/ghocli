/**
 * pages.go
 * Pages API
 *
 * Ghost Admin APIのPages機能を提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Page はGhostのページを表します
type Page struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug,omitempty"`
	HTML        string     `json:"html,omitempty"`
	Lexical     string     `json:"lexical,omitempty"`
	Status      string     `json:"status"` // draft, published, scheduled
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// PageListResponse はページ一覧のレスポンスです
type PageListResponse struct {
	Pages []Page `json:"pages"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// ListPages はページ一覧を取得します
func (c *Client) ListPages(opts ListOptions) (*PageListResponse, error) {
	path := "/ghost/api/admin/pages/"

	// クエリパラメータを構築
	params := []string{}
	if opts.Status != "" && opts.Status != "all" {
		params = append(params, fmt.Sprintf("filter=status:%s", opts.Status))
	}
	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Page > 0 {
		params = append(params, fmt.Sprintf("page=%d", opts.Page))
	}
	if opts.Include != "" {
		params = append(params, fmt.Sprintf("include=%s", opts.Include))
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
	var response PageListResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗: %w", err)
	}

	return &response, nil
}

// GetPage はページを取得します（IDまたはスラッグで指定）
func (c *Client) GetPage(idOrSlug string) (*Page, error) {
	// スラッグかどうかを判定（IDは通常24文字の16進数）
	var path string
	if len(idOrSlug) == 24 {
		// IDとして扱う
		path = fmt.Sprintf("/ghost/api/admin/pages/%s/", idOrSlug)
	} else {
		// スラッグとして扱う
		path = fmt.Sprintf("/ghost/api/admin/pages/slug/%s/", idOrSlug)
	}

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("ページが見つかりません: %s", idOrSlug)
	}

	return &response.Pages[0], nil
}

// CreatePage は新しいページを作成します
func (c *Client) CreatePage(page *Page) (*Page, error) {
	path := "/ghost/api/admin/pages/"

	// リクエストボディを作成
	reqBody := map[string]interface{}{
		"pages": []interface{}{page},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディの作成に失敗: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("POST", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("ページの作成に失敗しました")
	}

	return &response.Pages[0], nil
}

// UpdatePage はページを更新します
func (c *Client) UpdatePage(id string, page *Page) (*Page, error) {
	path := fmt.Sprintf("/ghost/api/admin/pages/%s/", id)

	// リクエストボディを作成
	reqBody := map[string]interface{}{
		"pages": []interface{}{page},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディの作成に失敗: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("ページの更新に失敗しました")
	}

	return &response.Pages[0], nil
}

// DeletePage はページを削除します
func (c *Client) DeletePage(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/pages/%s/", id)

	// リクエストを実行
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
