/**
 * users.go
 * Users API
 *
 * Ghost Admin APIのUsers機能を提供します。
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User はGhostのユーザーを表します
type User struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Slug         string    `json:"slug,omitempty"`
	Email        string    `json:"email,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Location     string    `json:"location,omitempty"`
	Website      string    `json:"website,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	CoverImage   string    `json:"cover_image,omitempty"`
	Roles        []Role    `json:"roles,omitempty"` // 読み取り専用
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Role はユーザーのロールを表します
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserListOptions はユーザー一覧取得のオプションです
type UserListOptions struct {
	Limit   int    // 取得件数（デフォルト: 15）
	Page    int    // ページ番号（デフォルト: 1）
	Include string // 含める追加情報（roles, count.posts など）
	Filter  string // フィルター条件
}

// UserListResponse はユーザー一覧のレスポンスです
type UserListResponse struct {
	Users []User `json:"users"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// UserResponse はユーザー単体のレスポンスです
type UserResponse struct {
	Users []User `json:"users"`
}

// ListUsers はユーザー一覧を取得します
func (c *Client) ListUsers(opts UserListOptions) (*UserListResponse, error) {
	path := "/ghost/api/admin/users/"

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
	var resp UserListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	return &resp, nil
}

// GetUser は指定されたIDまたはスラッグのユーザーを取得します
// idOrSlugが "slug:" で始まる場合はスラッグとして扱います
func (c *Client) GetUser(idOrSlug string) (*User, error) {
	var path string

	// スラッグかIDかを判定
	if strings.HasPrefix(idOrSlug, "slug:") {
		slug := strings.TrimPrefix(idOrSlug, "slug:")
		path = fmt.Sprintf("/ghost/api/admin/users/slug/%s/", slug)
	} else {
		path = fmt.Sprintf("/ghost/api/admin/users/%s/", idOrSlug)
	}

	// リクエストを実行
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp UserResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("ユーザーが見つかりません: %s", idOrSlug)
	}

	return &resp.Users[0], nil
}

// UpdateUser は既存のユーザーを更新します
func (c *Client) UpdateUser(id string, user *User) (*User, error) {
	path := fmt.Sprintf("/ghost/api/admin/users/%s/", id)

	// リクエストボディを構築
	reqBody := map[string]interface{}{
		"users": []User{*user},
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
	var resp UserResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("ユーザーの更新に失敗しました")
	}

	return &resp.Users[0], nil
}
