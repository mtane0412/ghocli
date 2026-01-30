/**
 * client_test.go
 * Ghost Admin APIクライアントのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewClient_クライアントの作成
func TestNewClient_クライアントの作成(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	siteURL := "https://test.ghost.io"

	client, err := NewClient(siteURL, keyID, secret)
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	if client == nil {
		t.Fatal("クライアントがnilです")
	}

	if client.baseURL != siteURL {
		t.Errorf("baseURL = %q; want %q", client.baseURL, siteURL)
	}
}

// TestNewClient_無効なURLでエラー
func TestNewClient_無効なURLでエラー(t *testing.T) {
	_, err := NewClient("", "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		t.Error("空のURLでエラーが返されなかった")
	}
}

// TestGetSite_サイト情報の取得
func TestGetSite_サイト情報の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/site/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/site/")
		}

		// Authorization ヘッダーが存在することを確認
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorizationヘッダーが設定されていない")
		}
		if len(auth) < 6 || auth[:6] != "Ghost " {
			t.Errorf("Authorizationヘッダーが不正: %s", auth)
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"site": map[string]interface{}{
				"title":       "Test Blog",
				"description": "A test blog",
				"url":         "https://test.ghost.io",
				"version":     "5.0",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// サイト情報を取得
	site, err := client.GetSite()
	if err != nil {
		t.Fatalf("サイト情報の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if site.Title != "Test Blog" {
		t.Errorf("Title = %q; want %q", site.Title, "Test Blog")
	}
	if site.Description != "A test blog" {
		t.Errorf("Description = %q; want %q", site.Description, "A test blog")
	}
	if site.URL != "https://test.ghost.io" {
		t.Errorf("URL = %q; want %q", site.URL, "https://test.ghost.io")
	}
	if site.Version != "5.0" {
		t.Errorf("Version = %q; want %q", site.Version, "5.0")
	}
}

// TestGetSite_APIエラー
func TestGetSite_APIエラー(t *testing.T) {
	// エラーを返すHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"message": "Unauthorized",
					"type":    "UnauthorizedError",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "invalid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// サイト情報を取得（エラーが返されることを期待）
	_, err = client.GetSite()
	if err == nil {
		t.Error("エラーが返されるべきだが、nilが返された")
	}
}
