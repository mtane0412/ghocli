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

// TestDetectFormat_Markdown はMarkdownファイルのフォーマット検出をテストする
func TestDetectFormat_Markdown(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".md拡張子",
			filePath: "test.md",
			expected: FormatMarkdown,
		},
		{
			name:     ".markdown拡張子",
			filePath: "article.markdown",
			expected: FormatMarkdown,
		},
		{
			name:     "パス付き.md",
			filePath: "/path/to/file.md",
			expected: FormatMarkdown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_HTML はHTMLファイルのフォーマット検出をテストする
func TestDetectFormat_HTML(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".html拡張子",
			filePath: "page.html",
			expected: FormatHTML,
		},
		{
			name:     ".htm拡張子",
			filePath: "index.htm",
			expected: FormatHTML,
		},
		{
			name:     "パス付き.html",
			filePath: "/var/www/page.html",
			expected: FormatHTML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_Lexical はLexical JSONファイルのフォーマット検出をテストする
func TestDetectFormat_Lexical(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".json拡張子",
			filePath: "content.json",
			expected: FormatLexical,
		},
		{
			name:     "パス付き.json",
			filePath: "/data/lexical.json",
			expected: FormatLexical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_Unknown は未知のファイル形式の検出をテストする
func TestDetectFormat_Unknown(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     "拡張子なし",
			filePath: "noextension",
			expected: FormatUnknown,
		},
		{
			name:     ".txt拡張子",
			filePath: "text.txt",
			expected: FormatUnknown,
		},
		{
			name:     "空文字列",
			filePath: "",
			expected: FormatUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestReadContentWithFormat_Markdown はMarkdownファイルの読み込みとフォーマット検出をテストする
func TestReadContentWithFormat_Markdown(t *testing.T) {
	// 一時Markdownファイルを作成
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	expectedContent := "# 見出し\n\nこれはMarkdownです。"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// ファイルから読み込み、フォーマットを検出
	content, format, err := ReadContentWithFormat(tmpFile, "")
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// コンテンツとフォーマットを検証
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	if format != FormatMarkdown {
		t.Errorf("format = %q; want %q", format, FormatMarkdown)
	}
}

// TestReadContentWithFormat_HTML はHTMLファイルの読み込みとフォーマット検出をテストする
func TestReadContentWithFormat_HTML(t *testing.T) {
	// 一時HTMLファイルを作成
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<h1>見出し</h1><p>これはHTMLです。</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// ファイルから読み込み、フォーマットを検出
	content, format, err := ReadContentWithFormat(tmpFile, "")
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// コンテンツとフォーマットを検証
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	if format != FormatHTML {
		t.Errorf("format = %q; want %q", format, FormatHTML)
	}
}

// TestReadContentWithFormat_InlineContent はインラインコンテンツの読み込みをテストする
func TestReadContentWithFormat_InlineContent(t *testing.T) {
	expectedContent := "<p>インラインコンテンツ</p>"

	// インラインコンテンツを読み込み
	content, format, err := ReadContentWithFormat("", expectedContent)
	if err != nil {
		t.Fatalf("コンテンツの読み込みに失敗: %v", err)
	}

	// コンテンツとフォーマットを検証
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	// インラインコンテンツの場合はフォーマット不明
	if format != FormatUnknown {
		t.Errorf("format = %q; want %q", format, FormatUnknown)
	}
}
