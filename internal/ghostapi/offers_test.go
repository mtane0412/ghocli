/**
 * offers_test.go
 * Offers APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListOffers_オファー一覧の取得
func TestListOffers_オファー一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/offers/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/offers/")
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
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "春のキャンペーン",
					"code":               "SPRING2024",
					"display_title":       "春の特別割引",
					"display_description": "今だけ50%オフ",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "プレミアム会員",
					},
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":                  "64fac5417c4c6b0001234568",
					"name":               "新規会員特典",
					"code":               "WELCOME100",
					"display_title":       "新規登録で100円オフ",
					"display_description": "初月限定",
					"type":               "fixed",
					"cadence":            "month",
					"amount":             100,
					"duration":           "once",
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   25,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "プレミアム会員",
					},
					"created_at": "2024-01-16T10:00:00.000Z",
					"updated_at": "2024-01-16T10:00:00.000Z",
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

	// オファー一覧を取得
	resp, err := client.ListOffers(OfferListOptions{})
	if err != nil {
		t.Fatalf("オファー一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Offers) != 2 {
		t.Errorf("オファー数 = %d; want 2", len(resp.Offers))
	}

	// 1つ目のオファーを検証
	firstOffer := resp.Offers[0]
	if firstOffer.Name != "春のキャンペーン" {
		t.Errorf("オファー名 = %q; want %q", firstOffer.Name, "春のキャンペーン")
	}
	if firstOffer.Code != "SPRING2024" {
		t.Errorf("コード = %q; want %q", firstOffer.Code, "SPRING2024")
	}
	if firstOffer.Type != "percent" {
		t.Errorf("タイプ = %q; want %q", firstOffer.Type, "percent")
	}
}

// TestListOffers_filterパラメータ
func TestListOffers_filterパラメータ(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		if !r.URL.Query().Has("filter") {
			t.Error("filterパラメータが設定されていない")
		}
		if r.URL.Query().Get("filter") != "status:active" {
			t.Errorf("filterパラメータ = %q; want %q", r.URL.Query().Get("filter"), "status:active")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "春のキャンペーン",
					"code":               "SPRING2024",
					"display_title":       "春の特別割引",
					"display_description": "今だけ50%オフ",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"created_at":         "2024-01-15T10:00:00.000Z",
					"updated_at":         "2024-01-15T10:00:00.000Z",
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

	// オファー一覧を取得（status:activeでフィルター）
	resp, err := client.ListOffers(OfferListOptions{
		Filter: "status:active",
	})
	if err != nil {
		t.Fatalf("オファー一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Offers) != 1 {
		t.Errorf("オファー数 = %d; want 1", len(resp.Offers))
	}
}

// TestGetOffer_IDでオファーを取得
func TestGetOffer_IDでオファーを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/offers/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"offers": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":               "春のキャンペーン",
					"code":               "SPRING2024",
					"display_title":       "春の特別割引",
					"display_description": "今だけ50%オフ",
					"type":               "percent",
					"cadence":            "month",
					"amount":             50,
					"duration":           "repeating",
					"duration_in_months": 3,
					"currency":           "JPY",
					"status":             "active",
					"redemption_count":   10,
					"tier": map[string]interface{}{
						"id":   "64fac5417c4c6b0001234999",
						"name": "プレミアム会員",
					},
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
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

	// オファーを取得
	offer, err := client.GetOffer("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("オファー取得エラー: %v", err)
	}

	// レスポンスの検証
	if offer.Name != "春のキャンペーン" {
		t.Errorf("オファー名 = %q; want %q", offer.Name, "春のキャンペーン")
	}
	if offer.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("オファーID = %q; want %q", offer.ID, "64fac5417c4c6b0001234567")
	}
}
