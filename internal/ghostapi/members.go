/**
 * members.go
 * Members API
 *
 * Ghost Admin APIのMembers機能を提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Member はGhostのメンバー（購読者）を表します
type Member struct {
	ID        string    `json:"id,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	Email     string    `json:"email"`              // 必須フィールド
	Name      string    `json:"name,omitempty"`
	Note      string    `json:"note,omitempty"`
	Status    string    `json:"status,omitempty"`   // free, paid, comped
	Labels    []Label   `json:"labels,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Label はメンバーに付与されるラベルを表します
type Label struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// MemberListOptions はメンバー一覧取得のオプションです
type MemberListOptions struct {
	Limit  int    // 取得件数（デフォルト: 15）
	Page   int    // ページ番号（デフォルト: 1）
	Filter string // フィルター条件
	Order  string // ソート順
}

// MemberListResponse はメンバー一覧のレスポンスです
type MemberListResponse struct {
	Members []Member `json:"members"`
	Meta    struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// MemberResponse はメンバー単体のレスポンスです
type MemberResponse struct {
	Members []Member `json:"members"`
}

// ListMembers はメンバー一覧を取得します
func (c *Client) ListMembers(opts MemberListOptions) (*MemberListResponse, error) {
	path := "/ghost/api/admin/members/"

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
	if opts.Order != "" {
		params = append(params, fmt.Sprintf("order=%s", opts.Order))
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
	var resp MemberListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetMember は指定されたIDのメンバーを取得します
func (c *Client) GetMember(id string) (*Member, error) {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("メンバーが見つかりません: %s", id)
	}

	return &resp.Members[0], nil
}

// CreateMember は新しいメンバーを作成します
func (c *Client) CreateMember(member *Member) (*Member, error) {
	path := "/ghost/api/admin/members/"

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"members": []Member{*member},
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
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("メンバーの作成に失敗しました")
	}

	return &resp.Members[0], nil
}

// UpdateMember は既存のメンバーを更新します
func (c *Client) UpdateMember(id string, member *Member) (*Member, error) {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"members": []Member{*member},
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
	var resp MemberResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("メンバーの更新に失敗しました")
	}

	return &resp.Members[0], nil
}

// DeleteMember はメンバーを削除します
func (c *Client) DeleteMember(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/members/%s/", id)

	// リクエストを実行
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return nil
}
