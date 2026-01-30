/**
 * users_test.go
 * Users APIのテスト
 *
 * Ghost Admin APIのUsers機能のテストを提供します。
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListUsers_ユーザー一覧の取得 は基本的なユーザー一覧取得をテストします
func TestListUsers_ユーザー一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドとパスを検証
		if r.Method != "GET" {
			t.Errorf("期待されるメソッド: GET, 実際: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/" {
			t.Errorf("期待されるパス: /ghost/api/admin/users/, 実際: %s", r.URL.Path)
		}

		// テストレスポンスを返す
		resp := UserListResponse{
			Users: []User{
				{
					ID:    "user1",
					Name:  "山田太郎",
					Slug:  "yamada-taro",
					Email: "yamada@example.com",
				},
				{
					ID:    "user2",
					Name:  "田中花子",
					Slug:  "tanaka-hanako",
					Email: "tanaka@example.com",
				},
			},
		}
		resp.Meta.Pagination.Page = 1
		resp.Meta.Pagination.Limit = 15
		resp.Meta.Pagination.Pages = 1
		resp.Meta.Pagination.Total = 2

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// ユーザー一覧を取得
	resp, err2 := client.ListUsers(UserListOptions{})
	if err2 != nil {
		t.Fatalf("ユーザー一覧の取得に失敗しました: %v", err2)
	}

	// レスポンスを検証
	if len(resp.Users) != 2 {
		t.Errorf("期待されるユーザー数: 2, 実際: %d", len(resp.Users))
	}

	// 1件目のユーザーを検証
	if resp.Users[0].Name != "山田太郎" {
		t.Errorf("期待される名前: 山田太郎, 実際: %s", resp.Users[0].Name)
	}
	if resp.Users[0].Email != "yamada@example.com" {
		t.Errorf("期待されるメールアドレス: yamada@example.com, 実際: %s", resp.Users[0].Email)
	}
}

// TestListUsers_オプション付き はクエリパラメータ付きのユーザー一覧取得をテストします
func TestListUsers_オプション付き(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータを検証
		query := r.URL.Query()
		if query.Get("limit") != "5" {
			t.Errorf("期待されるlimit: 5, 実際: %s", query.Get("limit"))
		}
		if query.Get("page") != "2" {
			t.Errorf("期待されるpage: 2, 実際: %s", query.Get("page"))
		}
		if query.Get("include") != "roles,count.posts" {
			t.Errorf("期待されるinclude: roles,count.posts, 実際: %s", query.Get("include"))
		}

		// テストレスポンスを返す
		resp := UserListResponse{
			Users: []User{
				{
					ID:    "user3",
					Name:  "佐藤一郎",
					Slug:  "sato-ichiro",
					Email: "sato@example.com",
					Roles: []Role{
						{ID: "role1", Name: "Author"},
					},
				},
			},
		}
		resp.Meta.Pagination.Page = 2
		resp.Meta.Pagination.Limit = 5
		resp.Meta.Pagination.Pages = 3
		resp.Meta.Pagination.Total = 15

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// オプション付きでユーザー一覧を取得
	opts := UserListOptions{
		Limit:   5,
		Page:    2,
		Include: "roles,count.posts",
	}
	resp, err2 := client.ListUsers(opts)
	if err2 != nil {
		t.Fatalf("ユーザー一覧の取得に失敗しました: %v", err2)
	}

	// レスポンスを検証
	if len(resp.Users) != 1 {
		t.Errorf("期待されるユーザー数: 1, 実際: %d", len(resp.Users))
	}
	if resp.Meta.Pagination.Page != 2 {
		t.Errorf("期待されるページ: 2, 実際: %d", resp.Meta.Pagination.Page)
	}
	if len(resp.Users[0].Roles) != 1 {
		t.Errorf("期待されるロール数: 1, 実際: %d", len(resp.Users[0].Roles))
	}
}

// TestGetUser_IDでユーザーを取得 はIDでユーザーを取得するテストです
func TestGetUser_IDでユーザーを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドとパスを検証
		if r.Method != "GET" {
			t.Errorf("期待されるメソッド: GET, 実際: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/user123/" {
			t.Errorf("期待されるパス: /ghost/api/admin/users/user123/, 実際: %s", r.URL.Path)
		}

		// テストレスポンスを返す
		resp := UserResponse{
			Users: []User{
				{
					ID:           "user123",
					Name:         "鈴木次郎",
					Slug:         "suzuki-jiro",
					Email:        "suzuki@example.com",
					Bio:          "エンジニア",
					Location:     "東京",
					Website:      "https://example.com",
					ProfileImage: "https://example.com/profile.jpg",
					CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// ユーザーを取得
	user, err2 := client.GetUser("user123")
	if err2 != nil {
		t.Fatalf("ユーザーの取得に失敗しました: %v", err2)
	}

	// レスポンスを検証
	if user.ID != "user123" {
		t.Errorf("期待されるID: user123, 実際: %s", user.ID)
	}
	if user.Name != "鈴木次郎" {
		t.Errorf("期待される名前: 鈴木次郎, 実際: %s", user.Name)
	}
	if user.Email != "suzuki@example.com" {
		t.Errorf("期待されるメールアドレス: suzuki@example.com, 実際: %s", user.Email)
	}
	if user.Bio != "エンジニア" {
		t.Errorf("期待される自己紹介: エンジニア, 実際: %s", user.Bio)
	}
}

// TestGetUser_スラッグでユーザーを取得 はスラッグでユーザーを取得するテストです
func TestGetUser_スラッグでユーザーを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドとパスを検証
		if r.Method != "GET" {
			t.Errorf("期待されるメソッド: GET, 実際: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/slug/suzuki-jiro/" {
			t.Errorf("期待されるパス: /ghost/api/admin/users/slug/suzuki-jiro/, 実際: %s", r.URL.Path)
		}

		// テストレスポンスを返す
		resp := UserResponse{
			Users: []User{
				{
					ID:    "user123",
					Name:  "鈴木次郎",
					Slug:  "suzuki-jiro",
					Email: "suzuki@example.com",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// スラッグでユーザーを取得
	user, err2 := client.GetUser("slug:suzuki-jiro")
	if err2 != nil {
		t.Fatalf("ユーザーの取得に失敗しました: %v", err2)
	}

	// レスポンスを検証
	if user.Slug != "suzuki-jiro" {
		t.Errorf("期待されるスラッグ: suzuki-jiro, 実際: %s", user.Slug)
	}
	if user.Name != "鈴木次郎" {
		t.Errorf("期待される名前: 鈴木次郎, 実際: %s", user.Name)
	}
}

// TestUpdateUser_ユーザーの更新 はユーザー更新をテストします
func TestUpdateUser_ユーザーの更新(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドとパスを検証
		if r.Method != "PUT" {
			t.Errorf("期待されるメソッド: PUT, 実際: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/user123/" {
			t.Errorf("期待されるパス: /ghost/api/admin/users/user123/, 実際: %s", r.URL.Path)
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースに失敗しました: %v", err)
		}

		users, ok := reqBody["users"].([]interface{})
		if !ok || len(users) == 0 {
			t.Fatal("リクエストボディにusersが含まれていません")
		}

		user := users[0].(map[string]interface{})
		if user["name"] != "更新後の名前" {
			t.Errorf("期待される名前: 更新後の名前, 実際: %v", user["name"])
		}

		// テストレスポンスを返す
		resp := UserResponse{
			Users: []User{
				{
					ID:       "user123",
					Name:     "更新後の名前",
					Slug:     "updated-slug",
					Email:    "updated@example.com",
					Bio:      "更新後の自己紹介",
					Location: "大阪",
					Website:  "https://updated.example.com",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// ユーザーを更新
	updateData := &User{
		Name:     "更新後の名前",
		Slug:     "updated-slug",
		Bio:      "更新後の自己紹介",
		Location: "大阪",
		Website:  "https://updated.example.com",
	}
	user, err2 := client.UpdateUser("user123", updateData)
	if err2 != nil {
		t.Fatalf("ユーザーの更新に失敗しました: %v", err2)
	}

	// レスポンスを検証
	if user.Name != "更新後の名前" {
		t.Errorf("期待される名前: 更新後の名前, 実際: %s", user.Name)
	}
	if user.Bio != "更新後の自己紹介" {
		t.Errorf("期待される自己紹介: 更新後の自己紹介, 実際: %s", user.Bio)
	}
	if user.Location != "大阪" {
		t.Errorf("期待される場所: 大阪, 実際: %s", user.Location)
	}
}

// TestGetUser_ユーザーが見つからない はユーザーが見つからない場合のエラー処理をテストします
func TestGetUser_ユーザーが見つからない(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 空のレスポンスを返す
		resp := UserResponse{
			Users: []User{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// ユーザーを取得（エラーが返ることを期待）
	_, err2 := client.GetUser("nonexistent")
	if err2 == nil {
		t.Fatal("エラーが期待されましたが、エラーが返されませんでした")
	}
}

// TestListUsers_APIエラー はAPIエラー時の処理をテストします
func TestListUsers_APIエラー(t *testing.T) {
	// テスト用のHTTPサーバーを作成（エラーを返す）
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":[{"message":"Internal Server Error"}]}`))
	}))
	defer ts.Close()

	// テスト用のクライアントを作成
	client, err := NewClient(ts.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアントの作成に失敗しました: %v", err)
	}

	// ユーザー一覧を取得（エラーが返ることを期待）
	_, err2 := client.ListUsers(UserListOptions{})
	if err2 == nil {
		t.Fatal("エラーが期待されましたが、エラーが返されませんでした")
	}
}
