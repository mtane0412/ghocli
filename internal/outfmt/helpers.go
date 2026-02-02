/**
 * helpers.go
 * 出力フォーマット用ヘルパー関数
 *
 * posts/pagesコマンドのデフォルト表示で使用する共通ヘルパー関数を提供します。
 */

package outfmt

import (
	"strings"

	"github.com/mtane0412/ghocli/internal/ghostapi"
)

// FormatAuthors は著者一覧を名前のカンマ区切りにフォーマットします
func FormatAuthors(authors []ghostapi.Author) string {
	// 著者がいない場合は空文字列
	if len(authors) == 0 {
		return ""
	}

	// 著者名を収集
	names := make([]string, len(authors))
	for i, author := range authors {
		names[i] = author.Name
	}

	// カンマ区切りで結合
	return strings.Join(names, ", ")
}

// FormatTags はタグ一覧を名前のカンマ区切りにフォーマットします
func FormatTags(tags []ghostapi.Tag) string {
	// タグがない場合は空文字列
	if len(tags) == 0 {
		return ""
	}

	// タグ名を収集
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}

	// カンマ区切りで結合
	return strings.Join(names, ", ")
}

// TruncateExcerpt は抜粋を指定文字数で切り詰めます
// maxLenを超える場合は、maxLen文字まで切り詰めて「...」を追加します
func TruncateExcerpt(excerpt string, maxLen int) string {
	// 空文字列の場合はそのまま返す
	if excerpt == "" {
		return ""
	}

	// 文字列をルーン（Unicodeコードポイント）のスライスに変換
	// 日本語などのマルチバイト文字を正しくカウントするため
	runes := []rune(excerpt)

	// maxLen以下の場合はそのまま返す
	if len(runes) <= maxLen {
		return excerpt
	}

	// maxLen文字まで切り詰めて「...」を追加
	return string(runes[:maxLen]) + "..."
}
