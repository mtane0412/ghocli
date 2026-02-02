/**
 * content_test.go
 * コンテンツ入力ユーティリティのテストコード
 */

package input

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadContent_ファイルから読み込み
func TestReadContent_ファイルから読み込み(t *testing.T) {
	// 一時ファイルを作成
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<p>テストコンテンツ</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// ファイルからコンテンツを読み込み
	content, err := ReadContent(tmpFile, "")
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// コンテンツの検証
	if content != expectedContent {
		t.Errorf("コンテンツ = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_インラインコンテンツを返す
func TestReadContent_インラインコンテンツを返す(t *testing.T) {
	expectedContent := "<p>インラインコンテンツ</p>"

	// インラインコンテンツを読み込み（ファイルパスが空）
	content, err := ReadContent("", expectedContent)
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// コンテンツの検証
	if content != expectedContent {
		t.Errorf("コンテンツ = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_ファイルが優先される
func TestReadContent_ファイルが優先される(t *testing.T) {
	// 一時ファイルを作成
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<p>ファイルコンテンツ</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// ファイルとインラインコンテンツの両方を指定（ファイルが優先される）
	content, err := ReadContent(tmpFile, "<p>インラインコンテンツ</p>")
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// ファイルのコンテンツが返されることを確認
	if content != expectedContent {
		t.Errorf("コンテンツ = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_ファイルが存在しない場合はエラー
func TestReadContent_ファイルが存在しない場合はエラー(t *testing.T) {
	// 存在しないファイルを指定
	_, err := ReadContent("/path/to/nonexistent/file.html", "")
	if err == nil {
		t.Error("エラーが返されるべき")
	}
}

// TestReadContent_両方空の場合は空文字列を返す
func TestReadContent_両方空の場合は空文字列を返す(t *testing.T) {
	// ファイルパスとインラインコンテンツの両方が空
	content, err := ReadContent("", "")
	if err != nil {
		t.Fatalf("エラーが返されるべきでない: %v", err)
	}

	// 空文字列が返されることを確認
	if content != "" {
		t.Errorf("コンテンツ = %q; want %q", content, "")
	}
}
