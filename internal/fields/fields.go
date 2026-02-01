/**
 * fields.go
 * フィールド定義基盤
 *
 * フィールドセットの定義、パース、バリデーション機能を提供します。
 */

package fields

import (
	"fmt"
	"strings"
)

// FieldSet はリソースのフィールドセットを表します
type FieldSet struct {
	// Default はlist用のデフォルトフィールド
	Default []string
	// Detail はget用のデフォルトフィールド
	Detail []string
	// All は全フィールド
	All []string
}

// Parse はカンマ区切りのフィールド指定文字列をパースします
//
// 入力例:
//   - "id,title,status" -> []string{"id", "title", "status"}
//   - "all" -> fieldSet.All
//   - "" -> nil（デフォルトフィールドを使用することを示す）
func Parse(input string, fieldSet FieldSet) ([]string, error) {
	// 空文字列の場合はnilを返す（デフォルトフィールドを使用）
	if input == "" {
		return nil, nil
	}

	// "all"の場合は全フィールドを返す
	if input == "all" {
		result := make([]string, len(fieldSet.All))
		copy(result, fieldSet.All)
		return result, nil
	}

	// カンマ区切りでパース
	fields := strings.Split(input, ",")

	// 各フィールドをトリム
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	// バリデーション
	if err := Validate(fields, fieldSet.All); err != nil {
		return nil, err
	}

	return fields, nil
}

// Validate は指定されたフィールドが利用可能かどうかを検証します
func Validate(fields []string, available []string) error {
	// 利用可能なフィールドをマップに変換
	availableMap := make(map[string]bool)
	for _, field := range available {
		availableMap[field] = true
	}

	// 各フィールドが利用可能かチェック
	for _, field := range fields {
		if !availableMap[field] {
			return fmt.Errorf("unknown field '%s'. Available fields: %s", field, strings.Join(available, ", "))
		}
	}

	return nil
}

// ListAvailable は利用可能なフィールド一覧を文字列として返します
func ListAvailable(fieldSet FieldSet) string {
	return "Specify fields with --fields. Available fields: " + strings.Join(fieldSet.All, ", ")
}
