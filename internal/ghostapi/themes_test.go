/**
 * themes_test.go
 * Themes APIのテストコード
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListThemes_テーマ一覧の取得
func TestListThemes_テーマ一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/themes/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/themes/")
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// Authorization ヘッダーが存在することを確認
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorizationヘッダーが設定されていない")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"themes": []map[string]interface{}{
				{
					"name":   "casper",
					"active": true,
					"package": map[string]interface{}{
						"name":        "casper",
						"description": "The default theme for Ghost",
						"version":     "5.0.0",
					},
					"templates": []map[string]interface{}{
						{"filename": "index.hbs"},
						{"filename": "post.hbs"},
					},
				},
				{
					"name":   "starter",
					"active": false,
					"package": map[string]interface{}{
						"name":        "starter",
						"description": "A minimal starter theme",
						"version":     "1.0.0",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// テーマ一覧を取得
	resp, err := client.ListThemes()
	if err != nil {
		t.Fatalf("テーマ一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Themes) != 2 {
		t.Errorf("テーマ数 = %d; want 2", len(resp.Themes))
	}

	// 1つ目のテーマを検証
	firstTheme := resp.Themes[0]
	if firstTheme.Name != "casper" {
		t.Errorf("テーマ名 = %q; want %q", firstTheme.Name, "casper")
	}
	if !firstTheme.Active {
		t.Error("Activeフラグ = false; want true")
	}
	if firstTheme.Package == nil {
		t.Fatal("Package情報がnil")
	}
	if firstTheme.Package.Name != "casper" {
		t.Errorf("Package名 = %q; want %q", firstTheme.Package.Name, "casper")
	}
	if len(firstTheme.Templates) != 2 {
		t.Errorf("テンプレート数 = %d; want 2", len(firstTheme.Templates))
	}
}

// TestUploadTheme_テーマのアップロード
func TestUploadTheme_テーマのアップロード(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/themes/upload/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/themes/upload/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// Content-Typeがmultipart/form-dataであることを確認
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || len(contentType) < 19 || contentType[:19] != "multipart/form-data" {
			t.Errorf("Content-Type = %q; want multipart/form-data", contentType)
		}

		// マルチパートフォームをパース
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("マルチパートフォームのパースエラー: %v", err)
		}

		// ファイルが存在することを確認
		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("ファイル取得エラー: %v", err)
		}
		if fileHeader.Filename != "theme.zip" {
			t.Errorf("ファイル名 = %q; want %q", fileHeader.Filename, "theme.zip")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"themes": []map[string]interface{}{
				{
					"name":   "custom-theme",
					"active": false,
					"package": map[string]interface{}{
						"name":        "custom-theme",
						"description": "A custom theme",
						"version":     "1.0.0",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// ダミーのZIPファイルデータを作成
	fileData := []byte("dummy zip content")
	reader := bytes.NewReader(fileData)

	// テーマをアップロード
	theme, err := client.UploadTheme(reader, "theme.zip")
	if err != nil {
		t.Fatalf("テーマアップロードエラー: %v", err)
	}

	// レスポンスの検証
	if theme.Name != "custom-theme" {
		t.Errorf("テーマ名 = %q; want %q", theme.Name, "custom-theme")
	}
	if theme.Active {
		t.Error("Activeフラグ = true; want false")
	}
	if theme.Package == nil {
		t.Fatal("Package情報がnil")
	}
	if theme.Package.Version != "1.0.0" {
		t.Errorf("バージョン = %q; want %q", theme.Package.Version, "1.0.0")
	}
}

// TestActivateTheme_テーマの有効化
func TestActivateTheme_テーマの有効化(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/themes/custom-theme/activate/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"themes": []map[string]interface{}{
				{
					"name":   "custom-theme",
					"active": true,
					"package": map[string]interface{}{
						"name":        "custom-theme",
						"description": "A custom theme",
						"version":     "1.0.0",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// テーマを有効化
	theme, err := client.ActivateTheme("custom-theme")
	if err != nil {
		t.Fatalf("テーマ有効化エラー: %v", err)
	}

	// レスポンスの検証
	if theme.Name != "custom-theme" {
		t.Errorf("テーマ名 = %q; want %q", theme.Name, "custom-theme")
	}
	if !theme.Active {
		t.Error("Activeフラグ = false; want true")
	}
}
