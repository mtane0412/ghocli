/**
 * client.go
 * Ghost Admin APIのHTTPクライアント
 *
 * Ghost Admin APIへのHTTPリクエストを管理します。
 * 各リクエストにはJWTトークンを含むAuthorizationヘッダーが付与されます。
 */

package ghostapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client はGhost Admin APIクライアントです
type Client struct {
	baseURL    string
	keyID      string
	secret     string
	httpClient *http.Client
}

// Site はGhostサイトの情報を表します
type Site struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Version     string `json:"version"`
}

// ErrorResponse はGhost APIのエラーレスポンスを表します
type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"errors"`
}

// NewClient は新しいGhost Admin APIクライアントを作成します。
func NewClient(baseURL, keyID, secret string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("サイトURLが空です")
	}
	if keyID == "" {
		return nil, errors.New("キーIDが空です")
	}
	if secret == "" {
		return nil, errors.New("シークレットが空です")
	}

	// URLの末尾のスラッシュを削除
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		keyID:   keyID,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// doRequest はHTTPリクエストを実行し、レスポンスボディを返します。
func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, error) {
	// JWTトークンを生成
	token, err := GenerateJWT(c.keyID, c.secret)
	if err != nil {
		return nil, fmt.Errorf("JWTの生成に失敗: %w", err)
	}

	// リクエストURLを構築
	url := c.baseURL + path

	// HTTPリクエストを作成
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("Authorization", "Ghost "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// リクエストを実行
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの実行に失敗: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込む
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗: %w", err)
	}

	// ステータスコードをチェック
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// エラーレスポンスをパース
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && len(errResp.Errors) > 0 {
			return nil, fmt.Errorf("APIエラー: %s", errResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	return respBody, nil
}

// GetSite はサイト情報を取得します。
func (c *Client) GetSite() (*Site, error) {
	respBody, err := c.doRequest("GET", "/ghost/api/admin/site/", nil)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var response struct {
		Site Site `json:"site"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗: %w", err)
	}

	return &response.Site, nil
}
