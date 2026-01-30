/**
 * newsletters_test.go
 * Newsletters APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListNewsletters_ニュースレター一覧の取得
func TestListNewsletters_ニュースレター一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/newsletters/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/newsletters/")
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
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "メインニュースレター",
					"description":         "週刊ニュースレター",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost編集部",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
				},
				{
					"id":                  "64fac5417c4c6b0001234568",
					"name":                "プレミアムニュースレター",
					"description":         "有料会員限定ニュースレター",
					"slug":                "premium-newsletter",
					"status":              "active",
					"visibility":          "paid",
					"subscribe_on_signup": false,
					"sender_name":         "Ghost編集部",
					"sender_email":        "premium@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          1,
					"created_at":          "2024-01-16T10:00:00.000Z",
					"updated_at":          "2024-01-16T10:00:00.000Z",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// ニュースレター一覧を取得
	resp, err := client.ListNewsletters(NewsletterListOptions{})
	if err != nil {
		t.Fatalf("ニュースレター一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Newsletters) != 2 {
		t.Errorf("ニュースレター数 = %d; want 2", len(resp.Newsletters))
	}

	// 1つ目のニュースレターを検証
	firstNewsletter := resp.Newsletters[0]
	if firstNewsletter.Name != "メインニュースレター" {
		t.Errorf("ニュースレター名 = %q; want %q", firstNewsletter.Name, "メインニュースレター")
	}
	if firstNewsletter.Slug != "main-newsletter" {
		t.Errorf("スラッグ = %q; want %q", firstNewsletter.Slug, "main-newsletter")
	}
	if firstNewsletter.Status != "active" {
		t.Errorf("ステータス = %q; want %q", firstNewsletter.Status, "active")
	}
}

// TestListNewsletters_filterパラメータ
func TestListNewsletters_filterパラメータ(t *testing.T) {
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
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "メインニュースレター",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// ニュースレター一覧を取得（status:activeでフィルター）
	resp, err := client.ListNewsletters(NewsletterListOptions{
		Filter: "status:active",
	})
	if err != nil {
		t.Fatalf("ニュースレター一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Newsletters) != 1 {
		t.Errorf("ニュースレター数 = %d; want 1", len(resp.Newsletters))
	}
}

// TestGetNewsletter_IDでニュースレターを取得
func TestGetNewsletter_IDでニュースレターを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/newsletters/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "メインニュースレター",
					"description":         "週刊ニュースレター",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost編集部",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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

	// ニュースレターを取得
	newsletter, err := client.GetNewsletter("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("ニュースレター取得エラー: %v", err)
	}

	// レスポンスの検証
	if newsletter.Name != "メインニュースレター" {
		t.Errorf("ニュースレター名 = %q; want %q", newsletter.Name, "メインニュースレター")
	}
	if newsletter.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("ニュースレターID = %q; want %q", newsletter.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetNewsletter_スラッグでニュースレターを取得
func TestGetNewsletter_スラッグでニュースレターを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/newsletters/slug/main-newsletter/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "メインニュースレター",
					"description":         "週刊ニュースレター",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "Ghost編集部",
					"sender_email":        "newsletter@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-15T10:00:00.000Z",
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

	// ニュースレターを取得
	newsletter, err := client.GetNewsletter("slug:main-newsletter")
	if err != nil {
		t.Fatalf("ニュースレター取得エラー: %v", err)
	}

	// レスポンスの検証
	if newsletter.Slug != "main-newsletter" {
		t.Errorf("スラッグ = %q; want %q", newsletter.Slug, "main-newsletter")
	}
}

// TestCreateNewsletter_ニュースレターの作成
func TestCreateNewsletter_ニュースレターの作成(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/newsletters/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/newsletters/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースに失敗: %v", err)
		}

		newsletters, ok := reqBody["newsletters"].([]interface{})
		if !ok || len(newsletters) == 0 {
			t.Error("newslettersフィールドが正しくない")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234569",
					"name":                "新規ニュースレター",
					"description":         "テスト用ニュースレター",
					"slug":                "new-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "テスト編集部",
					"sender_email":        "test@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-20T10:00:00.000Z",
					"updated_at":          "2024-01-20T10:00:00.000Z",
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

	// ニュースレターを作成
	newNewsletter := &Newsletter{
		Name:              "新規ニュースレター",
		Description:       "テスト用ニュースレター",
		Visibility:        "members",
		SubscribeOnSignup: true,
		SenderName:        "テスト編集部",
		SenderEmail:       "test@example.com",
	}

	createdNewsletter, err := client.CreateNewsletter(newNewsletter)
	if err != nil {
		t.Fatalf("ニュースレター作成エラー: %v", err)
	}

	// レスポンスの検証
	if createdNewsletter.Name != "新規ニュースレター" {
		t.Errorf("ニュースレター名 = %q; want %q", createdNewsletter.Name, "新規ニュースレター")
	}
	if createdNewsletter.ID == "" {
		t.Error("ニュースレターIDが設定されていない")
	}
}

// TestUpdateNewsletter_ニュースレターの更新
func TestUpdateNewsletter_ニュースレターの更新(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/newsletters/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースに失敗: %v", err)
		}

		newsletters, ok := reqBody["newsletters"].([]interface{})
		if !ok || len(newsletters) == 0 {
			t.Error("newslettersフィールドが正しくない")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"newsletters": []map[string]interface{}{
				{
					"id":                  "64fac5417c4c6b0001234567",
					"name":                "更新されたニュースレター",
					"description":         "更新後の説明",
					"slug":                "main-newsletter",
					"status":              "active",
					"visibility":          "members",
					"subscribe_on_signup": true,
					"sender_name":         "更新後編集部",
					"sender_email":        "updated@example.com",
					"sender_reply_to":     "newsletter",
					"sort_order":          0,
					"created_at":          "2024-01-15T10:00:00.000Z",
					"updated_at":          "2024-01-20T10:00:00.000Z",
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

	// ニュースレターを更新
	updateNewsletter := &Newsletter{
		Name:        "更新されたニュースレター",
		Description: "更新後の説明",
		SenderName:  "更新後編集部",
		SenderEmail: "updated@example.com",
	}

	updatedNewsletter, err := client.UpdateNewsletter("64fac5417c4c6b0001234567", updateNewsletter)
	if err != nil {
		t.Fatalf("ニュースレター更新エラー: %v", err)
	}

	// レスポンスの検証
	if updatedNewsletter.Name != "更新されたニュースレター" {
		t.Errorf("ニュースレター名 = %q; want %q", updatedNewsletter.Name, "更新されたニュースレター")
	}
	if updatedNewsletter.Description != "更新後の説明" {
		t.Errorf("説明 = %q; want %q", updatedNewsletter.Description, "更新後の説明")
	}
}
