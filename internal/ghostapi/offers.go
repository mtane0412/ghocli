/**
 * offers.go
 * Offers API
 *
 * Ghost Admin APIのOffers機能を提供します。
 * Create/Update操作には確認機構が適用されます。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Offer はGhostのオファーを表します
type Offer struct {
	ID                 string    `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Code               string    `json:"code,omitempty"`
	DisplayTitle       string    `json:"display_title,omitempty"`
	DisplayDescription string    `json:"display_description,omitempty"`
	Type               string    `json:"type,omitempty"`     // percent, fixed
	Cadence            string    `json:"cadence,omitempty"`  // month, year
	Amount             int       `json:"amount,omitempty"`
	Duration           string    `json:"duration,omitempty"` // once, forever, repeating
	DurationInMonths   int       `json:"duration_in_months,omitempty"`
	Currency           string    `json:"currency,omitempty"`
	Status             string    `json:"status,omitempty"`   // active, archived
	RedemptionCount    int       `json:"redemption_count,omitempty"`
	Tier               OfferTier `json:"tier,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// OfferTier はオファーに関連するティア情報を表します
type OfferTier struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// OfferListOptions はオファー一覧取得のオプションです
type OfferListOptions struct {
	Limit  int    // 取得件数（デフォルト: 15）
	Page   int    // ページ番号（デフォルト: 1）
	Filter string // フィルター条件
}

// OfferListResponse はオファー一覧のレスポンスです
type OfferListResponse struct {
	Offers []Offer `json:"offers"`
	Meta   struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// OfferResponse はオファー単体のレスポンスです
type OfferResponse struct {
	Offers []Offer `json:"offers"`
}

// ListOffers はオファー一覧を取得します
func (c *Client) ListOffers(opts OfferListOptions) (*OfferListResponse, error) {
	path := "/ghost/api/admin/offers/"

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
	var resp OfferListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetOffer は指定されたIDのオファーを取得します
// 注意: Offers APIはスラッグによる取得をサポートしていません（IDのみ）
func (c *Client) GetOffer(id string) (*Offer, error) {
	path := fmt.Sprintf("/ghost/api/admin/offers/%s/", id)

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("オファーが見つかりません: %s", id)
	}

	return &resp.Offers[0], nil
}

// CreateOffer は新しいオファーを作成します
func (c *Client) CreateOffer(offer *Offer) (*Offer, error) {
	path := "/ghost/api/admin/offers/"

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"offers": []interface{}{offer},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディのJSON化に失敗しました: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("POST", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("オファーの作成に失敗しました")
	}

	return &resp.Offers[0], nil
}

// UpdateOffer は既存のオファーを更新します
func (c *Client) UpdateOffer(id string, offer *Offer) (*Offer, error) {
	path := fmt.Sprintf("/ghost/api/admin/offers/%s/", id)

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"offers": []interface{}{offer},
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストボディのJSON化に失敗しました: %w", err)
	}

	// リクエストを実行
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp OfferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Offers) == 0 {
		return nil, fmt.Errorf("オファーの更新に失敗しました")
	}

	return &resp.Offers[0], nil
}
