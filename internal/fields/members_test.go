/**
 * members_test.go
 * Member用フィールド定義のテストコード
 */

package fields

import "testing"

// TestMemberFields_全フィールド数 はMemberFieldsが期待する数のフィールドを持つことを確認します
func TestMemberFields_全フィールド数(t *testing.T) {
	// Member構造体のフィールド数は9個
	expectedCount := 9
	if len(MemberFields.All) != expectedCount {
		t.Errorf("MemberFields.Allのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(MemberFields.All))
	}
}

// TestMemberFields_基本フィールド はMemberFieldsが基本フィールドを含むことを確認します
func TestMemberFields_基本フィールド(t *testing.T) {
	// 基本フィールドが存在することを確認
	expectedFields := []string{"id", "email", "name", "status"}

	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("基本フィールド '%s' が見つかりません", expected)
		}
	}
}

// TestMemberFields_詳細フィールド はMemberFieldsが詳細フィールドを含むことを確認します
func TestMemberFields_詳細フィールド(t *testing.T) {
	// 詳細フィールドが存在することを確認
	expectedFields := []string{"uuid", "note", "labels", "created_at", "updated_at"}

	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("詳細フィールド '%s' が見つかりません", expected)
		}
	}
}

// TestMemberFields_デフォルトフィールド はMemberFields.Defaultが期待するフィールドを含むことを確認します
func TestMemberFields_デフォルトフィールド(t *testing.T) {
	// Defaultフィールドは5個
	expectedCount := 5
	if len(MemberFields.Default) != expectedCount {
		t.Errorf("MemberFields.Defaultのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(MemberFields.Default))
	}

	// 必須フィールドが含まれることを確認
	expectedFields := []string{"id", "email", "name"}
	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Defaultフィールド '%s' が見つかりません", expected)
		}
	}
}

// TestMemberFields_詳細表示フィールド はMemberFields.Detailが期待するフィールドを含むことを確認します
func TestMemberFields_詳細表示フィールド(t *testing.T) {
	// Detailフィールドは9個（全フィールド）
	expectedCount := 9
	if len(MemberFields.Detail) != expectedCount {
		t.Errorf("MemberFields.Detailのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(MemberFields.Detail))
	}

	// uuid、note、labelsが含まれることを確認
	expectedFields := []string{"uuid", "note", "labels"}
	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detailフィールド '%s' が見つかりません", expected)
		}
	}
}
