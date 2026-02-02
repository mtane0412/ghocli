/**
 * pages_test.go
 * Pages APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListPages_ページ一覧の取得
func TestListPages_ページ一覧の取得(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
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
			"pages": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234601",
					"title":      "テストページ1",
					"slug":       "test-page-1",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234602",
					"title":      "テストページ2",
					"slug":       "test-page-2",
					"status":     "draft",
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
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// ページ一覧を取得
	response, err := client.ListPages(ListOptions{})
	if err != nil {
		t.Fatalf("ページ一覧の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if len(response.Pages) != 2 {
		t.Errorf("ページ数 = %d; want %d", len(response.Pages), 2)
	}
	if response.Pages[0].Title != "テストページ1" {
		t.Errorf("ページ1のタイトル = %q; want %q", response.Pages[0].Title, "テストページ1")
	}
	if response.Pages[0].Status != "published" {
		t.Errorf("ページ1のステータス = %q; want %q", response.Pages[0].Status, "published")
	}
	if response.Pages[1].Title != "テストページ2" {
		t.Errorf("ページ2のタイトル = %q; want %q", response.Pages[1].Title, "テストページ2")
	}
	if response.Pages[1].Status != "draft" {
		t.Errorf("ページ2のステータス = %q; want %q", response.Pages[1].Status, "draft")
	}
}

// TestGetPage_IDでページを取得
func TestGetPage_IDでページを取得(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "テストページ",
					"slug":       "test-page",
					"html":       "<p>ページ本文</p>",
					"status":     "published",
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
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// ページを取得
	page, err := client.GetPage(pageID)
	if err != nil {
		t.Fatalf("ページの取得に失敗: %v", err)
	}

	// レスポンスの検証
	if page.ID != pageID {
		t.Errorf("ID = %q; want %q", page.ID, pageID)
	}
	if page.Title != "テストページ" {
		t.Errorf("Title = %q; want %q", page.Title, "テストページ")
	}
	if page.HTML != "<p>ページ本文</p>" {
		t.Errorf("HTML = %q; want %q", page.HTML, "<p>ページ本文</p>")
	}
}

// TestCreatePage_ページの作成
func TestCreatePage_ページの作成(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// リクエストボディの検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディの読み込みに失敗: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "新規ページ" {
			t.Errorf("Title = %q; want %q", page["title"], "新規ページ")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234603",
					"title":      "新規ページ",
					"slug":       "new-page",
					"status":     "draft",
					"created_at": time.Now().Format(time.RFC3339),
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// ページを作成
	newPage := &Page{
		Title:  "新規ページ",
		Status: "draft",
	}
	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		t.Fatalf("ページの作成に失敗: %v", err)
	}

	// レスポンスの検証
	if createdPage.Title != "新規ページ" {
		t.Errorf("Title = %q; want %q", createdPage.Title, "新規ページ")
	}
	if createdPage.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPage.Status, "draft")
	}
	if createdPage.ID == "" {
		t.Error("IDが空です")
	}
}

// TestUpdatePage_ページの更新
func TestUpdatePage_ページの更新(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// リクエストボディの検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディの読み込みに失敗: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "更新後のページタイトル" {
			t.Errorf("Title = %q; want %q", page["title"], "更新後のページタイトル")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "更新後のページタイトル",
					"slug":       "updated-page",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": time.Now().Format(time.RFC3339),
				},
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

	// ページを更新
	updatePage := &Page{
		Title:  "更新後のページタイトル",
		Status: "published",
	}
	updatedPage, err := client.UpdatePage(pageID, updatePage)
	if err != nil {
		t.Fatalf("ページの更新に失敗: %v", err)
	}

	// レスポンスの検証
	if updatedPage.ID != pageID {
		t.Errorf("ID = %q; want %q", updatedPage.ID, pageID)
	}
	if updatedPage.Title != "更新後のページタイトル" {
		t.Errorf("Title = %q; want %q", updatedPage.Title, "更新後のページタイトル")
	}
}

// TestDeletePage_ページの削除
func TestDeletePage_ページの削除(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "DELETE")
		}

		// レスポンスを返す（204 No Content）
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// ページを削除
	err = client.DeletePage(pageID)
	if err != nil {
		t.Fatalf("ページの削除に失敗: %v", err)
	}
}

// TestGetPage_拡張フィールドのパース
func TestGetPage_拡張フィールドのパース(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// レスポンスを返す（拡張フィールドを含む）
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "拡張フィールドテストページ",
					"slug":       "extended-fields-test",
					"html":       "<p>ページ本文</p>",
					"status":     "published",
					"url":        "https://example.com/extended-fields-test/",
					"excerpt":    "ページの抜粋です。",
					"visibility": "public",
					"featured":   true,
					"authors": []map[string]interface{}{
						{
							"id":   "author1",
							"name": "山田太郎",
						},
					},
					"tags": []map[string]interface{}{
						{
							"id":   "tag1",
							"name": "テスト",
						},
						{
							"id":   "tag2",
							"name": "サンプル",
						},
					},
					"created_at":   "2024-01-15T10:00:00.000Z",
					"updated_at":   "2024-01-15T10:00:00.000Z",
					"published_at": "2024-01-15T10:00:00.000Z",
				},
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

	// ページを取得
	page, err := client.GetPage(pageID)
	if err != nil {
		t.Fatalf("ページの取得に失敗: %v", err)
	}

	// 基本フィールドの検証
	if page.ID != pageID {
		t.Errorf("ID = %q; want %q", page.ID, pageID)
	}
	if page.Title != "拡張フィールドテストページ" {
		t.Errorf("Title = %q; want %q", page.Title, "拡張フィールドテストページ")
	}

	// 拡張フィールドの検証
	if page.URL != "https://example.com/extended-fields-test/" {
		t.Errorf("URL = %q; want %q", page.URL, "https://example.com/extended-fields-test/")
	}
	if page.Excerpt != "ページの抜粋です。" {
		t.Errorf("Excerpt = %q; want %q", page.Excerpt, "ページの抜粋です。")
	}
	if page.Visibility != "public" {
		t.Errorf("Visibility = %q; want %q", page.Visibility, "public")
	}
	if !page.Featured {
		t.Errorf("Featured = %v; want %v", page.Featured, true)
	}

	// Authors の検証
	if len(page.Authors) != 1 {
		t.Errorf("Authors数 = %d; want %d", len(page.Authors), 1)
	}
	if len(page.Authors) > 0 && page.Authors[0].Name != "山田太郎" {
		t.Errorf("Authors[0].Name = %q; want %q", page.Authors[0].Name, "山田太郎")
	}

	// Tags の検証
	if len(page.Tags) != 2 {
		t.Errorf("Tags数 = %d; want %d", len(page.Tags), 2)
	}
	if len(page.Tags) > 0 && page.Tags[0].Name != "テスト" {
		t.Errorf("Tags[0].Name = %q; want %q", page.Tags[0].Name, "テスト")
	}
	if len(page.Tags) > 1 && page.Tags[1].Name != "サンプル" {
		t.Errorf("Tags[1].Name = %q; want %q", page.Tags[1].Name, "サンプル")
	}
}

// TestCreatePageWithOptions_HTMLソース指定でページを作成
func TestCreatePageWithOptions_HTMLソース指定でページを作成(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// source=htmlパラメータが含まれているか確認
		if r.URL.Query().Get("source") != "html" {
			t.Errorf("sourceパラメータ = %q; want %q", r.URL.Query().Get("source"), "html")
		}

		// リクエストボディの検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディの読み込みに失敗: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "HTMLページ" {
			t.Errorf("Title = %q; want %q", page["title"], "HTMLページ")
		}
		if page["html"] != "<h1>見出し</h1><p>段落</p>" {
			t.Errorf("HTML = %q; want %q", page["html"], "<h1>見出し</h1><p>段落</p>")
		}

		// レスポンスを返す（Lexical形式に変換された状態）
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":      "64fac5417c4c6b0001234604",
					"title":   "HTMLページ",
					"slug":    "html-page",
					"status":  "draft",
					"lexical": `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"見出し","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"heading","version":1,"tag":"h1"},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"段落","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
					"created_at": time.Now().Format(time.RFC3339),
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// HTMLページを作成
	newPage := &Page{
		Title:  "HTMLページ",
		HTML:   "<h1>見出し</h1><p>段落</p>",
		Status: "draft",
	}
	opts := CreateOptions{Source: "html"}
	createdPage, err := client.CreatePageWithOptions(newPage, opts)
	if err != nil {
		t.Fatalf("ページの作成に失敗: %v", err)
	}

	// レスポンスの検証
	if createdPage.Title != "HTMLページ" {
		t.Errorf("Title = %q; want %q", createdPage.Title, "HTMLページ")
	}
	if createdPage.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPage.Status, "draft")
	}
	// Lexical形式に変換されていることを確認
	if createdPage.Lexical == "" {
		t.Error("Lexicalフィールドが空です")
	}
}

// TestUpdatePageWithOptions_HTMLソース指定でページを更新
func TestUpdatePageWithOptions_HTMLソース指定でページを更新(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "PUT")
		}

		// source=htmlパラメータが含まれているか確認
		if r.URL.Query().Get("source") != "html" {
			t.Errorf("sourceパラメータ = %q; want %q", r.URL.Query().Get("source"), "html")
		}

		// リクエストボディの検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディの読み込みに失敗: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["html"] != "<h1>更新後の見出し</h1>" {
			t.Errorf("HTML = %q; want %q", page["html"], "<h1>更新後の見出し</h1>")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":      pageID,
					"title":   "更新後のページ",
					"slug":    "updated-html-page",
					"status":  "published",
					"lexical": `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"更新後の見出し","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"heading","version":1,"tag":"h1"}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": time.Now().Format(time.RFC3339),
				},
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

	// ページを更新
	updatePage := &Page{
		Title:  "更新後のページ",
		HTML:   "<h1>更新後の見出し</h1>",
		Status: "published",
	}
	opts := CreateOptions{Source: "html"}
	updatedPage, err := client.UpdatePageWithOptions(pageID, updatePage, opts)
	if err != nil {
		t.Fatalf("ページの更新に失敗: %v", err)
	}

	// レスポンスの検証
	if updatedPage.ID != pageID {
		t.Errorf("ID = %q; want %q", updatedPage.ID, pageID)
	}
	if updatedPage.Title != "更新後のページ" {
		t.Errorf("Title = %q; want %q", updatedPage.Title, "更新後のページ")
	}
	// Lexical形式に変換されていることを確認
	if updatedPage.Lexical == "" {
		t.Error("Lexicalフィールドが空です")
	}
}
