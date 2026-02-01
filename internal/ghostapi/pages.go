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
	// 基本情報
	ID     string `json:"id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	Title  string `json:"title"`
	Slug   string `json:"slug,omitempty"`
	Status string `json:"status"` // draft, published, scheduled
	URL    string `json:"url,omitempty"`

	// コンテンツ
	HTML          string `json:"html,omitempty"`
	Lexical       string `json:"lexical,omitempty"`
	Excerpt       string `json:"excerpt,omitempty"`
	CustomExcerpt string `json:"custom_excerpt,omitempty"`

	// 画像
	FeatureImage        string `json:"feature_image,omitempty"`
	FeatureImageAlt     string `json:"feature_image_alt,omitempty"`
	FeatureImageCaption string `json:"feature_image_caption,omitempty"`
	OGImage             string `json:"og_image,omitempty"`
	TwitterImage        string `json:"twitter_image,omitempty"`

	// SEO
	MetaTitle          string `json:"meta_title,omitempty"`
	MetaDescription    string `json:"meta_description,omitempty"`
	OGTitle            string `json:"og_title,omitempty"`
	OGDescription      string `json:"og_description,omitempty"`
	TwitterTitle       string `json:"twitter_title,omitempty"`
	TwitterDescription string `json:"twitter_description,omitempty"`
	CanonicalURL       string `json:"canonical_url,omitempty"`

	// 日時
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	// 制御
	Visibility string `json:"visibility,omitempty"` // public, members, paid
	Featured   bool   `json:"featured,omitempty"`
	EmailOnly  bool   `json:"email_only,omitempty"`

	// カスタム
	CodeinjectionHead string `json:"codeinjection_head,omitempty"`
	CodeinjectionFoot string `json:"codeinjection_foot,omitempty"`
	CustomTemplate    string `json:"custom_template,omitempty"`

	// 関連
	Tags          []Tag    `json:"tags,omitempty"`
	Authors       []Author `json:"authors,omitempty"`
	PrimaryAuthor *Author  `json:"primary_author,omitempty"`
	PrimaryTag    *Tag     `json:"primary_tag,omitempty"`

	// その他
	CommentID   string `json:"comment_id,omitempty"`
	ReadingTime int    `json:"reading_time,omitempty"`

	// メール・ニュースレター
	EmailSegment           string `json:"email_segment,omitempty"`
	NewsletterID           string `json:"newsletter_id,omitempty"`
	SendEmailWhenPublished bool   `json:"send_email_when_published,omitempty"`
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
		// IDとして扱う（formats=html,lexicalでHTML/Lexical両方を取得）
		path = fmt.Sprintf("/ghost/api/admin/pages/%s/?formats=html,lexical", idOrSlug)
	} else {
		// スラッグとして扱う（formats=html,lexicalでHTML/Lexical両方を取得）
		path = fmt.Sprintf("/ghost/api/admin/pages/slug/%s/?formats=html,lexical", idOrSlug)
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
