/**
 * newsletters.go
 * Newsletters API
 *
 * Ghost Admin APIのNewsletters機能を提供します。
 * ビジネス設定の誤変更リスクを回避するため、読み取り操作（List, Get）のみ実装しています。
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Newsletter はGhostのニュースレターを表します
type Newsletter struct {
	ID                string    `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	Slug              string    `json:"slug,omitempty"`
	Status            string    `json:"status,omitempty"`            // active, archived
	Visibility        string    `json:"visibility,omitempty"`        // members, paid
	SubscribeOnSignup bool      `json:"subscribe_on_signup,omitempty"`
	SenderName        string    `json:"sender_name,omitempty"`
	SenderEmail       string    `json:"sender_email,omitempty"`
	SenderReplyTo     string    `json:"sender_reply_to,omitempty"`
	SortOrder         int       `json:"sort_order,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

// NewsletterListOptions はニュースレター一覧取得のオプションです
type NewsletterListOptions struct {
	Limit  int    // 取得件数（デフォルト: 15）
	Page   int    // ページ番号（デフォルト: 1）
	Filter string // フィルター条件
}

// NewsletterListResponse はニュースレター一覧のレスポンスです
type NewsletterListResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
	Meta        struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// NewsletterResponse はニュースレター単体のレスポンスです
type NewsletterResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
}

// ListNewsletters はニュースレター一覧を取得します
func (c *Client) ListNewsletters(opts NewsletterListOptions) (*NewsletterListResponse, error) {
	path := "/ghost/api/admin/newsletters/"

	// クエリパラメータを構築
	params := []string{}
	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Page > 0 {
		params = append(params, fmt.Sprintf("page=%d", opts.Page))
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
	var resp NewsletterListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetNewsletter は指定されたIDまたはスラッグのニュースレターを取得します
// idOrSlugが "slug:" で始まる場合はスラッグとして扱います
func (c *Client) GetNewsletter(idOrSlug string) (*Newsletter, error) {
	var path string

	// スラッグかIDかを判定
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/newsletters/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/newsletters/%s/", idOrSlug)
	}

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp NewsletterResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Newsletters) == 0 {
		return nil, fmt.Errorf("ニュースレターが見つかりません: %s", idOrSlug)
	}

	return &resp.Newsletters[0], nil
}
