/**
 * members_test.go
 * Members APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListMembers_メンバー一覧の取得
func TestListMembers_メンバー一覧の取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/members/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/members/")
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
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "yamada@example.co.jp",
					"name":       "山田太郎",
					"note":       "テストメンバー",
					"status":     "free",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234568",
					"uuid":       "xyz987-uvw654-rst321",
					"email":      "tanaka@example.co.jp",
					"name":       "田中花子",
					"note":       "",
					"status":     "paid",
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

	// メンバー一覧を取得
	resp, err := client.ListMembers(MemberListOptions{})
	if err != nil {
		t.Fatalf("メンバー一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Members) != 2 {
		t.Errorf("メンバー数 = %d; want 2", len(resp.Members))
	}

	// 1つ目のメンバーを検証
	firstMember := resp.Members[0]
	if firstMember.Email != "yamada@example.co.jp" {
		t.Errorf("Email = %q; want %q", firstMember.Email, "yamada@example.co.jp")
	}
	if firstMember.Name != "山田太郎" {
		t.Errorf("Name = %q; want %q", firstMember.Name, "山田太郎")
	}
	if firstMember.Status != "free" {
		t.Errorf("Status = %q; want %q", firstMember.Status, "free")
	}
}

// TestListMembers_filterパラメータ
func TestListMembers_filterパラメータ(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		if !r.URL.Query().Has("filter") {
			t.Error("filterパラメータが設定されていない")
		}
		if r.URL.Query().Get("filter") != "status:paid" {
			t.Errorf("filterパラメータ = %q; want %q", r.URL.Query().Get("filter"), "status:paid")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234568",
					"email":      "tanaka@example.co.jp",
					"name":       "田中花子",
					"status":     "paid",
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
	client, err := NewClient(server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// メンバー一覧を取得（status:paidでフィルタ）
	resp, err := client.ListMembers(MemberListOptions{
		Filter: "status:paid",
	})
	if err != nil {
		t.Fatalf("メンバー一覧取得エラー: %v", err)
	}

	// レスポンスの検証
	if len(resp.Members) != 1 {
		t.Errorf("メンバー数 = %d; want 1", len(resp.Members))
	}
	if resp.Members[0].Status != "paid" {
		t.Errorf("Status = %q; want %q", resp.Members[0].Status, "paid")
	}
}

// TestGetMember_IDでメンバーを取得
func TestGetMember_IDでメンバーを取得(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "GET")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "yamada@example.co.jp",
					"name":       "山田太郎",
					"note":       "テストメンバー",
					"status":     "free",
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

	// メンバーを取得
	member, err := client.GetMember("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("メンバー取得エラー: %v", err)
	}

	// レスポンスの検証
	if member.Email != "yamada@example.co.jp" {
		t.Errorf("Email = %q; want %q", member.Email, "yamada@example.co.jp")
	}
	if member.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("ID = %q; want %q", member.ID, "64fac5417c4c6b0001234567")
	}
}

// TestCreateMember_メンバーの作成
func TestCreateMember_メンバーの作成(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/members/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/members/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// リクエストボディを検証
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("リクエストボディのパースエラー: %v", err)
		}

		members, ok := reqBody["members"].([]interface{})
		if !ok || len(members) == 0 {
			t.Error("リクエストボディに members 配列が存在しない")
		}

		// レスポンスを返す
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234999",
					"uuid":       "new-uuid-123",
					"email":      "new@example.co.jp",
					"name":       "新規メンバー",
					"note":       "新しく作成されたメンバー",
					"status":     "free",
					"created_at": createdAt,
					"updated_at": createdAt,
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

	// メンバーを作成
	newMember := &Member{
		Email: "new@example.co.jp",
		Name:  "新規メンバー",
		Note:  "新しく作成されたメンバー",
	}

	createdMember, err := client.CreateMember(newMember)
	if err != nil {
		t.Fatalf("メンバー作成エラー: %v", err)
	}

	// レスポンスの検証
	if createdMember.Email != "new@example.co.jp" {
		t.Errorf("Email = %q; want %q", createdMember.Email, "new@example.co.jp")
	}
	if createdMember.ID == "" {
		t.Error("IDが空")
	}
}

// TestUpdateMember_メンバーの更新
func TestUpdateMember_メンバーの更新(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
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
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "yamada@example.co.jp",
					"name":       "更新後の名前",
					"note":       "更新されたメンバー",
					"status":     "free",
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
	client, err := NewClient(server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("クライアント作成エラー: %v", err)
	}

	// メンバーを更新
	updateMember := &Member{
		Name: "更新後の名前",
		Note: "更新されたメンバー",
	}

	updatedMember, err := client.UpdateMember("64fac5417c4c6b0001234567", updateMember)
	if err != nil {
		t.Fatalf("メンバー更新エラー: %v", err)
	}

	// レスポンスの検証
	if updatedMember.Name != "更新後の名前" {
		t.Errorf("Name = %q; want %q", updatedMember.Name, "更新後の名前")
	}
}

// TestDeleteMember_メンバーの削除
func TestDeleteMember_メンバーの削除(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
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

	// メンバーを削除
	err = client.DeleteMember("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("メンバー削除エラー: %v", err)
	}
}
