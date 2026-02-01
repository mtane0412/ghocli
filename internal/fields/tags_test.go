/**
 * tags_test.go
 * Tag用フィールド定義のテストコード
 */

package fields

import "testing"

// TestTagFields_全フィールド数 はTagFieldsが期待する数のフィールドを持つことを確認します
func TestTagFields_全フィールド数(t *testing.T) {
	// Tag構造体のフィールド数は7個
	expectedCount := 7
	if len(TagFields.All) != expectedCount {
		t.Errorf("TagFields.Allのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(TagFields.All))
	}
}

// TestTagFields_基本フィールド はTagFieldsが基本フィールドを含むことを確認します
func TestTagFields_基本フィールド(t *testing.T) {
	// 基本フィールドが存在することを確認
	expectedFields := []string{"id", "name", "slug", "visibility"}

	fieldMap := make(map[string]bool)
	for _, field := range TagFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("基本フィールド '%s' が見つかりません", expected)
		}
	}
}

// TestTagFields_詳細フィールド はTagFieldsが詳細フィールドを含むことを確認します
func TestTagFields_詳細フィールド(t *testing.T) {
	// 詳細フィールドが存在することを確認
	expectedFields := []string{"description", "created_at", "updated_at"}

	fieldMap := make(map[string]bool)
	for _, field := range TagFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("詳細フィールド '%s' が見つかりません", expected)
		}
	}
}

// TestTagFields_デフォルトフィールド はTagFields.Defaultが期待するフィールドを含むことを確認します
func TestTagFields_デフォルトフィールド(t *testing.T) {
	// Defaultフィールドは5個
	expectedCount := 5
	if len(TagFields.Default) != expectedCount {
		t.Errorf("TagFields.Defaultのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(TagFields.Default))
	}

	// 必須フィールドが含まれることを確認
	expectedFields := []string{"id", "name", "slug"}
	fieldMap := make(map[string]bool)
	for _, field := range TagFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Defaultフィールド '%s' が見つかりません", expected)
		}
	}
}

// TestTagFields_詳細表示フィールド はTagFields.Detailが期待するフィールドを含むことを確認します
func TestTagFields_詳細表示フィールド(t *testing.T) {
	// Detailフィールドは7個（全フィールド）
	expectedCount := 7
	if len(TagFields.Detail) != expectedCount {
		t.Errorf("TagFields.Detailのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(TagFields.Detail))
	}

	// descriptionとupdated_atが含まれることを確認
	expectedFields := []string{"description", "updated_at"}
	fieldMap := make(map[string]bool)
	for _, field := range TagFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detailフィールド '%s' が見つかりません", expected)
		}
	}
}
