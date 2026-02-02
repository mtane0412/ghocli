/**
 * Markdown→HTML変換機能
 *
 * このパッケージはMarkdown形式のテキストをHTML形式に変換する機能を提供します。
 * goldmarkライブラリを使用して安全かつ高速な変換を実現します。
 */
package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// ConvertToHTML はMarkdown文字列をHTMLに変換する
//
// 引数:
//   - markdown: 変換元のMarkdown文字列
//
// 戻り値:
//   - string: 変換後のHTML文字列
//   - error: 変換エラー（通常はnilが返されます）
//
// 使用例:
//
//	html, err := ConvertToHTML("# 見出し\n\nこれは段落です。")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(html)
func ConvertToHTML(markdown string) (string, error) {
	// 空文字列の場合はそのまま返す
	if markdown == "" {
		return "", nil
	}

	// バッファを用意
	var buf bytes.Buffer

	// goldmarkを使ってMarkdown→HTML変換
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
