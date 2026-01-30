/**
 * tiers_test.go
 * Tiers APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListTiers_ティア一覧の取得
func TestListTiers_ティア一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/tiers/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/tiers/")
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
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "無料会員",
					"description":      "無料で記事を読めます",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
				},
				{
					"id":               "64fac5417c4c6b0001234568",
					"name":             "プレミアム会員",
					"description":      "すべての記事にアクセス可能",
					"slug":             "premium",
					"active":           true,
					"type":             "paid",
					"visibility":       "public",
					"monthly_price":    500,
					"yearly_price":     5000,
					"currency":         "JPY",
					"welcome_page_url": "/welcome",
					"created_at":       "2024-01-16T10:00:00.000Z",
					"updated_at":       "2024-01-16T10:00:00.000Z",
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 2,
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

	// ティア一覧を取得
	resp, err := client.ListTiers(TierListOptions{})
	if err != nil {
		t.Fatalf("ティア一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Tiers) != 2 {
		t.Errorf("ティア数 = %d; want 2", len(resp.Tiers))
	}

	// 1つ目のティアを検証
	firstTier := resp.Tiers[0]
	if firstTier.Name != "無料会員" {
		t.Errorf("ティア名 = %q; want %q", firstTier.Name, "無料会員")
	}
	if firstTier.Slug != "free" {
		t.Errorf("スラッグ = %q; want %q", firstTier.Slug, "free")
	}
	if firstTier.Type != "free" {
		t.Errorf("タイプ = %q; want %q", firstTier.Type, "free")
	}
}

// TestListTiers_includeパラメータ
func TestListTiers_includeパラメータ(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		if !r.URL.Query().Has("include") {
			t.Error("includeパラメータが設定されていない")
		}
		if r.URL.Query().Get("include") != "monthly_price,yearly_price" {
			t.Errorf("includeパラメータ = %q; want %q", r.URL.Query().Get("include"), "monthly_price,yearly_price")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":            "64fac5417c4c6b0001234568",
					"name":          "プレミアム会員",
					"slug":          "premium",
					"type":          "paid",
					"active":        true,
					"visibility":    "public",
					"monthly_price": 500,
					"yearly_price":  5000,
					"currency":      "JPY",
					"created_at":    "2024-01-16T10:00:00.000Z",
					"updated_at":    "2024-01-16T10:00:00.000Z",
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 1,
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

	// ティア一覧を取得（monthly_price, yearly_priceを含む）
	resp, err := client.ListTiers(TierListOptions{
		Include: "monthly_price,yearly_price",
	})
	if err != nil {
		t.Fatalf("ティア一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Tiers) != 1 {
		t.Errorf("ティア数 = %d; want 1", len(resp.Tiers))
	}
}

// TestGetTier_IDでティアを取得
func TestGetTier_IDでティアを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tiers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "無料会員",
					"description":      "無料で記事を読めます",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
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

	// ティアを取得
	tier, err := client.GetTier("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("ティア取得エラー: %v", err)
	}

	// レスポンスの検証
	if tier.Name != "無料会員" {
		t.Errorf("ティア名 = %q; want %q", tier.Name, "無料会員")
	}
	if tier.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("ティアID = %q; want %q", tier.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetTier_スラッグでティアを取得
func TestGetTier_スラッグでティアを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tiers/slug/free/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tiers": []map[string]interface{}{
				{
					"id":               "64fac5417c4c6b0001234567",
					"name":             "無料会員",
					"description":      "無料で記事を読めます",
					"slug":             "free",
					"active":           true,
					"type":             "free",
					"visibility":       "public",
					"welcome_page_url": "",
					"created_at":       "2024-01-15T10:00:00.000Z",
					"updated_at":       "2024-01-15T10:00:00.000Z",
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

	// ティアを取得
	tier, err := client.GetTier("slug:free")
	if err != nil {
		t.Fatalf("ティア取得エラー: %v", err)
	}

	// レスポンスの検証
	if tier.Slug != "free" {
		t.Errorf("スラッグ = %q; want %q", tier.Slug, "free")
	}
}
