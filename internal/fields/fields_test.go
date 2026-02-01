/**
 * fields_test.go
 * フィールド定義のテスト
 *
 * フィールドパーサーとバリデーション機能のテストを提供します。
 */

package fields

import (
	"strings"
	"testing"
)

// TestParse_正常系_カンマ区切り指定 はカンマ区切りフィールド指定をパースできることを確認します
func TestParse_正常系_カンマ区切り指定(t *testing.T) {
	// テストデータ：利用可能なフィールド
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status", "html"},
		All:     []string{"id", "title", "status", "html", "slug", "url"},
	}

	// フィールド指定をパース
	result, err := Parse("id,title,status", fieldSet)
	if err != nil {
		t.Fatalf("パースに失敗: %v", err)
	}

	// 期待値
	expected := []string{"id", "title", "status"}

	// 結果を検証
	if len(result) != len(expected) {
		t.Errorf("結果の長さが不正: got=%d, want=%d", len(result), len(expected))
	}
	for i, field := range expected {
		if result[i] != field {
			t.Errorf("結果[%d]が不正: got=%s, want=%s", i, result[i], field)
		}
	}
}

// TestParse_正常系_all指定 は"all"指定で全フィールドを取得できることを確認します
func TestParse_正常系_all指定(t *testing.T) {
	// テストデータ
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status", "html", "slug"},
	}

	// "all"を指定
	result, err := Parse("all", fieldSet)
	if err != nil {
		t.Fatalf("パースに失敗: %v", err)
	}

	// 全フィールドが返されることを確認
	if len(result) != len(fieldSet.All) {
		t.Errorf("結果の長さが不正: got=%d, want=%d", len(result), len(fieldSet.All))
	}
	for i, field := range fieldSet.All {
		if result[i] != field {
			t.Errorf("結果[%d]が不正: got=%s, want=%s", i, result[i], field)
		}
	}
}

// TestParse_異常系_無効なフィールド は無効なフィールド指定でエラーを返すことを確認します
func TestParse_異常系_無効なフィールド(t *testing.T) {
	// テストデータ
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status"},
	}

	// 無効なフィールドを指定
	_, err := Parse("id,invalid_field", fieldSet)
	if err == nil {
		t.Fatal("エラーが返されませんでした")
	}

	// エラーメッセージに"invalid_field"が含まれることを確認
	if !strings.Contains(err.Error(), "invalid_field") {
		t.Errorf("エラーメッセージに無効なフィールド名が含まれていません: %v", err)
	}
}

// TestValidate_正常系 は有効なフィールドリストでエラーが返されないことを確認します
func TestValidate_正常系(t *testing.T) {
	// テストデータ
	available := []string{"id", "title", "status", "html"}
	fields := []string{"id", "title"}

	// バリデーション
	err := Validate(fields, available)
	if err != nil {
		t.Errorf("バリデーションエラー: %v", err)
	}
}

// TestValidate_異常系_無効なフィールド は無効なフィールドでエラーを返すことを確認します
func TestValidate_異常系_無効なフィールド(t *testing.T) {
	// テストデータ
	available := []string{"id", "title", "status"}
	fields := []string{"id", "invalid"}

	// バリデーション
	err := Validate(fields, available)
	if err == nil {
		t.Fatal("エラーが返されませんでした")
	}

	// エラーメッセージに"invalid"が含まれることを確認
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("エラーメッセージに無効なフィールド名が含まれていません: %v", err)
	}
}

// TestListAvailable はフィールド一覧を文字列として取得できることを確認します
func TestListAvailable(t *testing.T) {
	// テストデータ
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status", "html", "slug"},
	}

	// フィールド一覧を取得
	result := ListAvailable(fieldSet)

	// 結果が空でないことを確認
	if result == "" {
		t.Error("フィールド一覧が空です")
	}

	// 全フィールドが含まれることを確認
	for _, field := range fieldSet.All {
		if !strings.Contains(result, field) {
			t.Errorf("フィールド一覧に'%s'が含まれていません: %s", field, result)
		}
	}
}
