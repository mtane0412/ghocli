/**
 * posts_test.go
 * Posts APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListPosts_投稿一覧の取得
func TestListPosts_投稿一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
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
			"posts": []map[string]interface{}{
				{
					"id":           "64fac5417c4c6b0001234567",
					"title":        "テスト投稿1",
					"slug":         "test-post-1",
					"status":       "published",
					"created_at":   "2024-01-15T10:00:00.000Z",
					"updated_at":   "2024-01-15T10:00:00.000Z",
					"published_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234568",
					"title":      "テスト投稿2",
					"slug":       "test-post-2",
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿一覧を取得
	response, err := client.ListPosts(ListOptions{})
	if err != nil {
		t.Fatalf("投稿一覧の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if len(response.Posts) != 2 {
		t.Errorf("投稿数 = %d; want %d", len(response.Posts), 2)
	}
	if response.Posts[0].Title != "テスト投稿1" {
		t.Errorf("投稿1のタイトル = %q; want %q", response.Posts[0].Title, "テスト投稿1")
	}
	if response.Posts[0].Status != "published" {
		t.Errorf("投稿1のステータス = %q; want %q", response.Posts[0].Status, "published")
	}
	if response.Posts[1].Title != "テスト投稿2" {
		t.Errorf("投稿2のタイトル = %q; want %q", response.Posts[1].Title, "テスト投稿2")
	}
	if response.Posts[1].Status != "draft" {
		t.Errorf("投稿2のステータス = %q; want %q", response.Posts[1].Status, "draft")
	}
}

// TestListPosts_ステータスフィルタ
func TestListPosts_ステータスフィルタ(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		status := r.URL.Query().Get("filter")
		if status != "status:draft" {
			t.Errorf("ステータスフィルタ = %q; want %q", status, "status:draft")
		}

		// レスポンスを返す（draftのみ）
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234568",
					"title":      "下書き投稿",
					"slug":       "draft-post",
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
					"total": 1,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// draftステータスでフィルタリング
	response, err := client.ListPosts(ListOptions{Status: "draft"})
	if err != nil {
		t.Fatalf("投稿一覧の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if len(response.Posts) != 1 {
		t.Errorf("投稿数 = %d; want %d", len(response.Posts), 1)
	}
	if response.Posts[0].Status != "draft" {
		t.Errorf("ステータス = %q; want %q", response.Posts[0].Status, "draft")
	}
}

// TestGetPost_IDで投稿を取得
func TestGetPost_IDで投稿を取得(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":           postID,
					"title":        "テスト投稿",
					"slug":         "test-post",
					"html":         "<p>本文</p>",
					"status":       "published",
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿を取得
	post, err := client.GetPost(postID)
	if err != nil {
		t.Fatalf("投稿の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if post.ID != postID {
		t.Errorf("ID = %q; want %q", post.ID, postID)
	}
	if post.Title != "テスト投稿" {
		t.Errorf("Title = %q; want %q", post.Title, "テスト投稿")
	}
	if post.HTML != "<p>本文</p>" {
		t.Errorf("HTML = %q; want %q", post.HTML, "<p>本文</p>")
	}
}

// TestGetPost_スラッグで投稿を取得
func TestGetPost_スラッグで投稿を取得(t *testing.T) {
	slug := "test-post"

	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/posts/slug/" + slug + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":           "64fac5417c4c6b0001234567",
					"title":        "テスト投稿",
					"slug":         slug,
					"html":         "<p>本文</p>",
					"status":       "published",
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿を取得
	post, err := client.GetPost(slug)
	if err != nil {
		t.Fatalf("投稿の取得に失敗: %v", err)
	}

	// レスポンスの検証
	if post.Slug != slug {
		t.Errorf("Slug = %q; want %q", post.Slug, slug)
	}
	if post.Title != "テスト投稿" {
		t.Errorf("Title = %q; want %q", post.Title, "テスト投稿")
	}
}

// TestCreatePost_投稿の作成
func TestCreatePost_投稿の作成(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// リクエストボディの検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディの読み込みに失敗: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["title"] != "新規投稿" {
			t.Errorf("Title = %q; want %q", post["title"], "新規投稿")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234569",
					"title":      "新規投稿",
					"slug":       "new-post",
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿を作成
	newPost := &Post{
		Title:  "新規投稿",
		Status: "draft",
	}
	createdPost, err := client.CreatePost(newPost)
	if err != nil {
		t.Fatalf("投稿の作成に失敗: %v", err)
	}

	// レスポンスの検証
	if createdPost.Title != "新規投稿" {
		t.Errorf("Title = %q; want %q", createdPost.Title, "新規投稿")
	}
	if createdPost.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPost.Status, "draft")
	}
	if createdPost.ID == "" {
		t.Error("IDが空です")
	}
}

// TestUpdatePost_投稿の更新
func TestUpdatePost_投稿の更新(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
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
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["title"] != "更新後のタイトル" {
			t.Errorf("Title = %q; want %q", post["title"], "更新後のタイトル")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         postID,
					"title":      "更新後のタイトル",
					"slug":       "updated-post",
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿を更新
	updatePost := &Post{
		Title:  "更新後のタイトル",
		Status: "published",
	}
	updatedPost, err := client.UpdatePost(postID, updatePost)
	if err != nil {
		t.Fatalf("投稿の更新に失敗: %v", err)
	}

	// レスポンスの検証
	if updatedPost.ID != postID {
		t.Errorf("ID = %q; want %q", updatedPost.ID, postID)
	}
	if updatedPost.Title != "更新後のタイトル" {
		t.Errorf("Title = %q; want %q", updatedPost.Title, "更新後のタイトル")
	}
}

// TestDeletePost_投稿の削除
func TestDeletePost_投稿の削除(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
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
	client, err := NewClient(server.URL, "keyid", "secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗: %v", err)
	}

	// 投稿を削除
	err = client.DeletePost(postID)
	if err != nil {
		t.Fatalf("投稿の削除に失敗: %v", err)
	}
}
