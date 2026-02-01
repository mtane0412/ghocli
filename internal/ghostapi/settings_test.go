/**
 * settings_test.go
 * Settings APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetSettings_設定一覧の取得
func TestGetSettings_設定一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/settings/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/settings/")
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
			"settings": []map[string]interface{}{
				{
					"key":   "title",
					"value": "My Ghost Site",
				},
				{
					"key":   "description",
					"value": "Thoughts, stories and ideas",
				},
				{
					"key":   "timezone",
					"value": "Asia/Tokyo",
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

	// 設定一覧を取得
	resp, err := client.GetSettings()
	if err != nil {
		t.Fatalf("設定一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Settings) != 3 {
		t.Errorf("設定数 = %d; want 3", len(resp.Settings))
	}

	// 1つ目の設定を検証
	firstSetting := resp.Settings[0]
	if firstSetting.Key != "title" {
		t.Errorf("設定キー = %q; want %q", firstSetting.Key, "title")
	}
	if firstSetting.Value != "My Ghost Site" {
		t.Errorf("設定値 = %q; want %q", firstSetting.Value, "My Ghost Site")
	}
}

// TestUpdateSettings_設定の更新
func TestUpdateSettings_設定の更新(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/settings/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/settings/")
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// リクエストボディをパース
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースエラー: %v", err)
		}

		// settingsが含まれていることを確認
		settings, ok := reqBody["settings"].([]interface{})
		if !ok {
			t.Fatal("settingsフィールドが存在しない")
		}

		if len(settings) != 1 {
			t.Errorf("設定数 = %d; want 1", len(settings))
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"settings": []map[string]interface{}{
				{
					"key":   "title",
					"value": "Updated Title",
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

	// 設定を更新
	updates := []SettingUpdate{
		{
			Key:   "title",
			Value: "Updated Title",
		},
	}
	resp, err := client.UpdateSettings(updates)
	if err != nil {
		t.Fatalf("設定更新エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Settings) != 1 {
		t.Errorf("設定数 = %d; want 1", len(resp.Settings))
	}

	updatedSetting := resp.Settings[0]
	if updatedSetting.Key != "title" {
		t.Errorf("設定キー = %q; want %q", updatedSetting.Key, "title")
	}
	if updatedSetting.Value != "Updated Title" {
		t.Errorf("設定値 = %q; want %q", updatedSetting.Value, "Updated Title")
	}
}
