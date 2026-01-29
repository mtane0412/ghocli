/**
 * images_test.go
 * Images APIのテストコード
 */

package ghostapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUploadImage_画像のアップロード
func TestUploadImage_画像のアップロード(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.URL.Path != "/ghost/api/admin/images/upload/" {
			t.Errorf("リクエストパス = %q; want %q", r.URL.Path, "/ghost/api/admin/images/upload/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTPメソッド = %q; want %q", r.Method, "POST")
		}

		// Authorization ヘッダーが存在することを確認
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorizationヘッダーが設定されていない")
		}

		// Content-Typeがmultipart/form-dataであることを確認
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Errorf("Content-Type = %q; want multipart/form-data", contentType)
		}

		// マルチパートフォームをパース
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			t.Fatalf("マルチパートフォームのパースエラー: %v", err)
		}

		// ファイルが存在することを確認
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("ファイルの取得エラー: %v", err)
		}
		defer file.Close()

		// ファイル名を確認
		if header.Filename != "test-image.jpg" {
			t.Errorf("ファイル名 = %q; want %q", header.Filename, "test-image.jpg")
		}

		// ファイル内容を確認
		fileContent, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("ファイル内容の読み込みエラー: %v", err)
		}
		expectedContent := "fake image content"
		if string(fileContent) != expectedContent {
			t.Errorf("ファイル内容 = %q; want %q", string(fileContent), expectedContent)
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"images": []map[string]interface{}{
				{
					"url": "https://example.com/content/images/2024/01/test-image.jpg",
					"ref": nil,
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

	// 画像をアップロード
	fakeImageContent := strings.NewReader("fake image content")
	image, err := client.UploadImage(fakeImageContent, "test-image.jpg", ImageUploadOptions{})
	if err != nil {
		t.Fatalf("画像アップロードエラー: %v", err)
	}

	// レスポンスの検証
	expectedURL := "https://example.com/content/images/2024/01/test-image.jpg"
	if image.URL != expectedURL {
		t.Errorf("画像URL = %q; want %q", image.URL, expectedURL)
	}
}

// TestUploadImage_purposeパラメータ
func TestUploadImage_purposeパラメータ(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// マルチパートフォームをパース
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			t.Fatalf("マルチパートフォームのパースエラー: %v", err)
		}

		// purposeパラメータを確認
		purpose := r.FormValue("purpose")
		if purpose != "profile_image" {
			t.Errorf("purpose = %q; want %q", purpose, "profile_image")
		}

		// refパラメータを確認
		ref := r.FormValue("ref")
		if ref != "test-ref-12345" {
			t.Errorf("ref = %q; want %q", ref, "test-ref-12345")
		}

		// レスポンスを返す
		response := map[string]interface{}{
			"images": []map[string]interface{}{
				{
					"url": "https://example.com/content/images/2024/01/profile.jpg",
					"ref": "test-ref-12345",
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

	// 画像をアップロード（purposeとrefを指定）
	fakeImageContent := strings.NewReader("fake profile image")
	image, err := client.UploadImage(fakeImageContent, "profile.jpg", ImageUploadOptions{
		Purpose: "profile_image",
		Ref:     "test-ref-12345",
	})
	if err != nil {
		t.Fatalf("画像アップロードエラー: %v", err)
	}

	// レスポンスの検証
	if image.Ref != "test-ref-12345" {
		t.Errorf("ref = %q; want %q", image.Ref, "test-ref-12345")
	}
}
