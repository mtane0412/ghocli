/**
 * Markdown→HTML変換機能のテストコード
 */
package markdown

import (
	"strings"
	"testing"
)

// TestConvertToHTML_Paragraph は段落の変換をテストする
func TestConvertToHTML_Paragraph(t *testing.T) {
	markdown := "これは段落です。"
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("変換エラーが発生しました: %v", err)
	}

	// goldmarkは<p>タグで囲み、末尾に改行を追加する
	expected := "<p>これは段落です。</p>\n"
	if html != expected {
		t.Errorf("期待値と異なります。\n期待値: %q\n実際値: %q", expected, html)
	}
}

// TestConvertToHTML_Headings は見出し（h1-h6）の変換をテストする
func TestConvertToHTML_Headings(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "h1見出し",
			markdown: "# 見出し1",
			expected: "<h1>見出し1</h1>\n",
		},
		{
			name:     "h2見出し",
			markdown: "## 見出し2",
			expected: "<h2>見出し2</h2>\n",
		},
		{
			name:     "h3見出し",
			markdown: "### 見出し3",
			expected: "<h3>見出し3</h3>\n",
		},
		{
			name:     "h6見出し",
			markdown: "###### 見出し6",
			expected: "<h6>見出し6</h6>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("変換エラーが発生しました: %v", err)
			}
			if html != tt.expected {
				t.Errorf("期待値と異なります。\n期待値: %q\n実際値: %q", tt.expected, html)
			}
		})
	}
}

// TestConvertToHTML_Lists はリスト（箇条書き、番号付き）の変換をテストする
func TestConvertToHTML_Lists(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string // 完全一致ではなく、含まれるべき要素をチェック
	}{
		{
			name:     "箇条書きリスト",
			markdown: "- アイテム1\n- アイテム2\n- アイテム3",
			contains: []string{"<ul>", "<li>アイテム1</li>", "<li>アイテム2</li>", "<li>アイテム3</li>", "</ul>"},
		},
		{
			name:     "番号付きリスト",
			markdown: "1. 最初\n2. 次\n3. 最後",
			contains: []string{"<ol>", "<li>最初</li>", "<li>次</li>", "<li>最後</li>", "</ol>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("変換エラーが発生しました: %v", err)
			}

			// 各要素が含まれているか確認
			for _, elem := range tt.contains {
				if !strings.Contains(html, elem) {
					t.Errorf("期待される要素が含まれていません: %q\nHTML: %q", elem, html)
				}
			}
		})
	}
}

// TestConvertToHTML_CodeBlock はコードブロックの変換をテストする
func TestConvertToHTML_CodeBlock(t *testing.T) {
	markdown := "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("変換エラーが発生しました: %v", err)
	}

	// コードブロックは<pre><code>で囲まれる
	if !strings.Contains(html, "<pre>") {
		t.Errorf("<pre>タグが含まれていません: %q", html)
	}
	if !strings.Contains(html, "<code") {
		t.Errorf("<code>タグが含まれていません: %q", html)
	}
	if !strings.Contains(html, "func main()") {
		t.Errorf("コード内容が含まれていません: %q", html)
	}
}

// TestConvertToHTML_LinksAndImages はリンクと画像の変換をテストする
func TestConvertToHTML_LinksAndImages(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "リンク",
			markdown: "[Googleへのリンク](https://google.com)",
			contains: []string{"<a href=\"https://google.com\"", "Googleへのリンク</a>"},
		},
		{
			name:     "画像",
			markdown: "![代替テキスト](/path/to/image.png)",
			contains: []string{"<img src=\"/path/to/image.png\"", "alt=\"代替テキスト\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("変換エラーが発生しました: %v", err)
			}

			for _, elem := range tt.contains {
				if !strings.Contains(html, elem) {
					t.Errorf("期待される要素が含まれていません: %q\nHTML: %q", elem, html)
				}
			}
		})
	}
}

// TestConvertToHTML_EmptyString は空文字列入力時の動作をテストする
func TestConvertToHTML_EmptyString(t *testing.T) {
	markdown := ""
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("空文字列でエラーが発生しました: %v", err)
	}

	// 空文字列の場合も空文字列を返すべき
	if html != "" {
		t.Errorf("空文字列を期待しましたが、実際は: %q", html)
	}
}

// TestConvertToHTML_ComplexMarkdown は複雑なMarkdownの変換をテストする
func TestConvertToHTML_ComplexMarkdown(t *testing.T) {
	markdown := `# 見出し

これは**太字**と*斜体*を含む段落です。

- リスト項目1
- リスト項目2

[リンク](https://example.com)も含みます。`

	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("変換エラーが発生しました: %v", err)
	}

	// 各要素が含まれているか確認
	expectedElements := []string{
		"<h1>見出し</h1>",
		"<strong>太字</strong>",
		"<em>斜体</em>",
		"<ul>",
		"<li>リスト項目1</li>",
		"<a href=\"https://example.com\"",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(html, elem) {
			t.Errorf("期待される要素が含まれていません: %q\nHTML: %q", elem, html)
		}
	}
}
