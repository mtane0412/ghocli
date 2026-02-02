/**
 * helpers_test.go
 * 出力フォーマット用ヘルパー関数のテスト
 */

package outfmt

import (
	"testing"

	"github.com/mtane0412/ghocli/internal/ghostapi"
)

// TestFormatAuthors_著者一覧のフォーマット
func TestFormatAuthors_著者一覧のフォーマット(t *testing.T) {
	// テストケース: 複数の著者
	authors := []ghostapi.Author{
		{ID: "1", Name: "山田太郎"},
		{ID: "2", Name: "鈴木花子"},
	}

	// フォーマットを実行
	result := FormatAuthors(authors)

	// 期待値を検証
	expected := "山田太郎, 鈴木花子"
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatAuthors_単一の著者
func TestFormatAuthors_単一の著者(t *testing.T) {
	// テストケース: 単一の著者
	authors := []ghostapi.Author{
		{ID: "1", Name: "山田太郎"},
	}

	// フォーマットを実行
	result := FormatAuthors(authors)

	// 期待値を検証
	expected := "山田太郎"
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatAuthors_著者なし
func TestFormatAuthors_著者なし(t *testing.T) {
	// テストケース: 空のスライス
	authors := []ghostapi.Author{}

	// フォーマットを実行
	result := FormatAuthors(authors)

	// 期待値を検証（空文字列）
	expected := ""
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatTags_タグ一覧のフォーマット
func TestFormatTags_タグ一覧のフォーマット(t *testing.T) {
	// テストケース: 複数のタグ
	tags := []ghostapi.Tag{
		{ID: "1", Name: "旅行"},
		{ID: "2", Name: "北海道"},
		{ID: "3", Name: "グルメ"},
	}

	// フォーマットを実行
	result := FormatTags(tags)

	// 期待値を検証
	expected := "旅行, 北海道, グルメ"
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestFormatTags_単一のタグ
func TestFormatTags_単一のタグ(t *testing.T) {
	// テストケース: 単一のタグ
	tags := []ghostapi.Tag{
		{ID: "1", Name: "旅行"},
	}

	// フォーマットを実行
	result := FormatTags(tags)

	// 期待値を検証
	expected := "旅行"
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestFormatTags_タグなし
func TestFormatTags_タグなし(t *testing.T) {
	// テストケース: 空のスライス
	tags := []ghostapi.Tag{}

	// フォーマットを実行
	result := FormatTags(tags)

	// 期待値を検証（空文字列）
	expected := ""
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestTruncateExcerpt_抜粋の切り詰め
func TestTruncateExcerpt_抜粋の切り詰め(t *testing.T) {
	// テストケース: 長い文字列（actualに140文字を超える文字列）
	excerpt := "これは非常に長い抜粋テキストです。この抜粋は140文字を超えるため、適切に切り詰められる必要があります。切り詰められた部分には「...」が追加されます。これは人間やLLMにとって読みやすくするための処理です。さらに文字を追加して140文字を超えるようにします。あと少しで140文字に達します。もう少し追加します。これで140文字を確実に超えるはずです。"

	// フォーマットを実行（最大140文字）
	result := TruncateExcerpt(excerpt, 140)

	// 期待値を検証
	// ルーン数（文字数）で検証
	resultRunes := []rune(result)
	excerptRunes := []rune(excerpt)

	// 元の文字列が140文字を超えているか確認
	if len(excerptRunes) <= 140 {
		t.Errorf("テストケースのexcerptが140文字以下です。len = %d", len(excerptRunes))
	}

	// 結果が143文字（140 + "..."）であることを確認
	if len(resultRunes) != 143 {
		t.Errorf("TruncateExcerpt()の長さ = %d; want %d", len(resultRunes), 143)
	}

	// 末尾が「...」であることを確認
	if len(resultRunes) >= 3 {
		suffix := string(resultRunes[len(resultRunes)-3:])
		if suffix != "..." {
			t.Errorf("TruncateExcerpt()の末尾 = %q; want %q", suffix, "...")
		}
	}
}

// TestTruncateExcerpt_短い文字列
func TestTruncateExcerpt_短い文字列(t *testing.T) {
	// テストケース: 短い文字列
	excerpt := "これは短い抜粋です。"

	// フォーマットを実行（最大140文字）
	result := TruncateExcerpt(excerpt, 140)

	// 期待値を検証（そのまま返される）
	if result != excerpt {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, excerpt)
	}
}

// TestTruncateExcerpt_ちょうど最大長
func TestTruncateExcerpt_ちょうど最大長(t *testing.T) {
	// テストケース: ちょうど140文字（日本語1文字 = 3バイトだが、文字数でカウント）
	excerpt := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"

	// フォーマットを実行（最大140文字）
	result := TruncateExcerpt(excerpt, 140)

	// 期待値を検証（そのまま返される）
	if result != excerpt {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, excerpt)
	}
}

// TestTruncateExcerpt_空文字列
func TestTruncateExcerpt_空文字列(t *testing.T) {
	// テストケース: 空文字列
	excerpt := ""

	// フォーマットを実行
	result := TruncateExcerpt(excerpt, 140)

	// 期待値を検証（空文字列）
	if result != "" {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, "")
	}
}
