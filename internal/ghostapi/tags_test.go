/**
 * tags_test.go
 * Tags APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListTags_タグ一覧の取得
func TestListTags_タグ一覧の取得(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/tags/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/tags/")
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
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "テクノロジー関連の記事",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
				},
				{
					"id":          "64fac5417c4c6b0001234568",
					"name":        "Programming",
					"slug":        "programming",
					"description": "プログラミングのTips",
					"visibility":  "public",
					"created_at":  "2024-01-16T10:00:00.000Z",
					"updated_at":  "2024-01-16T10:00:00.000Z",
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

	// タグ一覧を取得
	resp, err := client.ListTags(TagListOptions{})
	if err != nil {
		t.Fatalf("タグ一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Tags) != 2 {
		t.Errorf("タグ数 = %d; want 2", len(resp.Tags))
	}

	// 1つ目のタグを検証
	firstTag := resp.Tags[0]
	if firstTag.Name != "Technology" {
		t.Errorf("タグ名 = %q; want %q", firstTag.Name, "Technology")
	}
	if firstTag.Slug != "technology" {
		t.Errorf("スラッグ = %q; want %q", firstTag.Slug, "technology")
	}
}

// TestListTags_includeパラメータ
func TestListTags_includeパラメータ(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		if !r.URL.Query().Has("include") {
			t.Error("includeパラメータが設定されていない")
		}
		if r.URL.Query().Get("include") != "count.posts" {
			t.Errorf("includeパラメータ = %q; want %q", r.URL.Query().Get("include"), "count.posts")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
					"count": map[string]interface{}{
						"posts": 10,
					},
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

	// タグ一覧を取得（count.postsを含む）
	resp, err := client.ListTags(TagListOptions{
		Include: "count.posts",
	})
	if err != nil {
		t.Fatalf("タグ一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Tags) != 1 {
		t.Errorf("タグ数 = %d; want 1", len(resp.Tags))
	}
}

// TestGetTag_IDでタグを取得
func TestGetTag_IDでタグを取得(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "テクノロジー関連の記事",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
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

	// タグを取得
	tag, err := client.GetTag("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("タグ取得エラー: %v", err)
	}

	// レスポンスの検証
	if tag.Name != "Technology" {
		t.Errorf("タグ名 = %q; want %q", tag.Name, "Technology")
	}
	if tag.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("タグID = %q; want %q", tag.ID, "64fac5417c4c6b0001234567")
	}
}

// TestGetTag_スラッグでタグを取得
func TestGetTag_スラッグでタグを取得(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tags/slug/technology/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Technology",
					"slug":        "technology",
					"description": "テクノロジー関連の記事",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  "2024-01-15T10:00:00.000Z",
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

	// タグを取得
	tag, err := client.GetTag("slug:technology")
	if err != nil {
		t.Fatalf("タグ取得エラー: %v", err)
	}

	// レスポンスの検証
	if tag.Slug != "technology" {
		t.Errorf("スラッグ = %q; want %q", tag.Slug, "technology")
	}
}

// TestCreateTag_タグの作成
func TestCreateTag_タグの作成(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/tags/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/tags/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースエラー: %v", err)
		}

		tags, ok := reqBody["tags"].([]interface{})
		if !ok || len(tags) == 0 {
			t.Error("リクエストボディに tags 配列が存在しない")
		}

		// レスポンスを返す
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234999",
					"name":        "New Tag",
					"slug":        "new-tag",
					"description": "新しいタグ",
					"visibility":  "public",
					"created_at":  createdAt,
					"updated_at":  createdAt,
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

	// タグを作成
	newTag := &Tag{
		Name:        "New Tag",
		Description: "新しいタグ",
	}

	createdTag, err := client.CreateTag(newTag)
	if err != nil {
		t.Fatalf("タグ作成エラー: %v", err)
	}

	// レスポンスの検証
	if createdTag.Name != "New Tag" {
		t.Errorf("タグ名 = %q; want %q", createdTag.Name, "New Tag")
	}
	if createdTag.ID == "" {
		t.Error("タグIDが空")
	}
}

// TestUpdateTag_タグの更新
func TestUpdateTag_タグの更新(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
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
		response := map[string]interface{}{
			"tags": []map[string]interface{}{
				{
					"id":          "64fac5417c4c6b0001234567",
					"name":        "Updated Technology",
					"slug":        "technology",
					"description": "更新されたテクノロジー関連の記事",
					"visibility":  "public",
					"created_at":  "2024-01-15T10:00:00.000Z",
					"updated_at":  time.Now().Format(time.RFC3339),
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

	// タグを更新
	updateTag := &Tag{
		Name:        "Updated Technology",
		Description: "更新されたテクノロジー関連の記事",
	}

	updatedTag, err := client.UpdateTag("64fac5417c4c6b0001234567", updateTag)
	if err != nil {
		t.Fatalf("タグ更新エラー: %v", err)
	}

	// レスポンスの検証
	if updatedTag.Name != "Updated Technology" {
		t.Errorf("タグ名 = %q; want %q", updatedTag.Name, "Updated Technology")
	}
}

// TestDeleteTag_タグの削除
func TestDeleteTag_タグの削除(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/tags/64fac5417c4c6b0001234567/"
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// タグを削除
	err = client.DeleteTag("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("タグ削除エラー: %v", err)
	}
}
