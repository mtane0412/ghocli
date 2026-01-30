/**
 * webhooks_test.go
 * Webhooks APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreateWebhook_Webhookの作成
func TestCreateWebhook_Webhookの作成(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/webhooks/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/webhooks/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// Authorization ヘッダーが存在することを確認
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorizationヘッダーが設定されていない")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースエラー: %v", err)
		}

		webhooks, ok := reqBody["webhooks"].([]interface{})
		if !ok || len(webhooks) == 0 {
			t.Error("リクエストボディに webhooks 配列が存在しない")
		}

		// レスポンスを返す
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"webhooks": []map[string]interface{}{
				{
					"id":                "64fac5417c4c6b0001234567",
					"event":             "post.published",
					"target_url":        "https://example.com/webhook",
					"name":              "Post published webhook",
					"secret":            "secret123",
					"api_version":       "v5.0",
					"integration_id":    "64fac5417c4c6b0001234568",
					"status":            "available",
					"last_triggered_at": nil,
					"created_at":        createdAt,
					"updated_at":        createdAt,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// Webhookを作成
	webhook := &Webhook{
		Event:     "post.published",
		TargetURL: "https://example.com/webhook",
		Name:      "Post published webhook",
	}

	created, err := client.CreateWebhook(webhook)
	if err != nil {
		t.Fatalf("Webhook作成エラー: %v", err)
	}

	// レスポンスの検証
	if created.Event != "post.published" {
		t.Errorf("イベント = %q; want %q", created.Event, "post.published")
	}
	if created.TargetURL != "https://example.com/webhook" {
		t.Errorf("ターゲットURL = %q; want %q", created.TargetURL, "https://example.com/webhook")
	}
	if created.ID == "" {
		t.Error("Webhook IDが空")
	}
	if created.Status != "available" {
		t.Errorf("ステータス = %q; want %q", created.Status, "available")
	}
}

// TestUpdateWebhook_Webhookの更新
func TestUpdateWebhook_Webhookの更新(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/webhooks/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースエラー: %v", err)
		}

		// レスポンスを返す
		updatedAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"webhooks": []map[string]interface{}{
				{
					"id":                "64fac5417c4c6b0001234567",
					"event":             "post.published",
					"target_url":        "https://example.com/webhook-updated",
					"name":              "Updated webhook",
					"secret":            "secret123",
					"api_version":       "v5.0",
					"integration_id":    "64fac5417c4c6b0001234568",
					"status":            "available",
					"last_triggered_at": nil,
					"created_at":        "2024-01-15T10:00:00.000Z",
					"updated_at":        updatedAt,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// Webhookを更新
	updateWebhook := &Webhook{
		TargetURL: "https://example.com/webhook-updated",
		Name:      "Updated webhook",
	}

	updated, err := client.UpdateWebhook("64fac5417c4c6b0001234567", updateWebhook)
	if err != nil {
		t.Fatalf("Webhook更新エラー: %v", err)
	}

	// レスポンスの検証
	if updated.TargetURL != "https://example.com/webhook-updated" {
		t.Errorf("ターゲットURL = %q; want %q", updated.TargetURL, "https://example.com/webhook-updated")
	}
	if updated.Name != "Updated webhook" {
		t.Errorf("名前 = %q; want %q", updated.Name, "Updated webhook")
	}
}

// TestDeleteWebhook_Webhookの削除
func TestDeleteWebhook_Webhookの削除(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/webhooks/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "DELETE")
		}

		// 204 No Content を返す
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// Webhookを削除
	err = client.DeleteWebhook("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Webhook削除エラー: %v", err)
	}
}
