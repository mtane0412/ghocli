/**
 * users_test.go
 * User用フィールド定義のテストコード
 */

package fields

import "testing"

// TestUserFields_全フィールド数 はUserFieldsが期待する数のフィールドを持つことを確認します
func TestUserFields_全フィールド数(t *testing.T) {
	// User構造体のフィールド数は12個
	expectedCount := 12
	if len(UserFields.All) != expectedCount {
		t.Errorf("UserFields.Allのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(UserFields.All))
	}
}

// TestUserFields_基本フィールド はUserFieldsが基本フィールドを含むことを確認します
func TestUserFields_基本フィールド(t *testing.T) {
	// 基本フィールドが存在することを確認
	expectedFields := []string{"id", "name", "slug", "email"}

	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("基本フィールド '%s' が見つかりません", expected)
		}
	}
}

// TestUserFields_プロフィールフィールド はUserFieldsがプロフィールフィールドを含むことを確認します
func TestUserFields_プロフィールフィールド(t *testing.T) {
	// プロフィールフィールドが存在することを確認
	expectedFields := []string{"bio", "location", "website", "profile_image", "cover_image"}

	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("プロフィールフィールド '%s' が見つかりません", expected)
		}
	}
}

// TestUserFields_ロールフィールド はUserFieldsがrolesフィールドを含むことを確認します
func TestUserFields_ロールフィールド(t *testing.T) {
	// rolesフィールドが存在することを確認
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	if !fieldMap["roles"] {
		t.Errorf("rolesフィールドが見つかりません")
	}
}

// TestUserFields_デフォルトフィールド はUserFields.Defaultが期待するフィールドを含むことを確認します
func TestUserFields_デフォルトフィールド(t *testing.T) {
	// Defaultフィールドは5個
	expectedCount := 5
	if len(UserFields.Default) != expectedCount {
		t.Errorf("UserFields.Defaultのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(UserFields.Default))
	}

	// 必須フィールドが含まれることを確認
	expectedFields := []string{"id", "name", "email"}
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Defaultフィールド '%s' が見つかりません", expected)
		}
	}
}

// TestUserFields_詳細表示フィールド はUserFields.Detailが期待するフィールドを含むことを確認します
func TestUserFields_詳細表示フィールド(t *testing.T) {
	// Detailフィールドは12個（全フィールド）
	expectedCount := 12
	if len(UserFields.Detail) != expectedCount {
		t.Errorf("UserFields.Detailのフィールド数が正しくありません。expected=%d, got=%d", expectedCount, len(UserFields.Detail))
	}

	// bio、location、rolesが含まれることを確認
	expectedFields := []string{"bio", "location", "roles"}
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detailフィールド '%s' が見つかりません", expected)
		}
	}
}
