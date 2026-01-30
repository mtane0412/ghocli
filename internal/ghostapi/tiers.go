/**
 * tiers.go
 * Tiers API
 *
 * Ghost Admin APIのTiers機能を提供します。
 * ビジネス設定の誤変更リスクを回避するため、読み取り操作（List, Get）のみ実装しています。
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tier はGhostのティアを表します
type Tier struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	Active         bool      `json:"active,omitempty"`
	Type           string    `json:"type,omitempty"`          // free, paid
	Visibility     string    `json:"visibility,omitempty"`    // public, none
	WelcomePageURL string    `json:"welcome_page_url,omitempty"`
	MonthlyPrice   int       `json:"monthly_price,omitempty"` // 最小通貨単位
	YearlyPrice    int       `json:"yearly_price,omitempty"`
	Currency       string    `json:"currency,omitempty"`
	Benefits       []string  `json:"benefits,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// TierListOptions はティア一覧取得のオプションです
type TierListOptions struct {
	Limit   int    // 取得件数（デフォルト: 15）
	Page    int    // ページ番号（デフォルト: 1）
	Include string // 含める追加情報（monthly_price, yearly_price, benefits など）
	Filter  string // フィルター条件
}

// TierListResponse はティア一覧のレスポンスです
type TierListResponse struct {
	Tiers []Tier `json:"tiers"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// TierResponse はティア単体のレスポンスです
type TierResponse struct {
	Tiers []Tier `json:"tiers"`
}

// ListTiers はティア一覧を取得します
func (c *Client) ListTiers(opts TierListOptions) (*TierListResponse, error) {
	path := "/ghost/api/admin/tiers/"

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
	var resp TierListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetTier は指定されたIDまたはスラッグのティアを取得します
// idOrSlugが "slug:" で始まる場合はスラッグとして扱います
func (c *Client) GetTier(idOrSlug string) (*Tier, error) {
	var path string

	// スラッグかIDかを判定
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/tiers/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/tiers/%s/", idOrSlug)
	}

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp TierResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Tiers) == 0 {
		return nil, fmt.Errorf("ティアが見つかりません: %s", idOrSlug)
	}

	return &resp.Tiers[0], nil
}
