/**
 * confirm_test.go
 * 破壊的操作の確認機構のテストコード
 */

package cmd

import (
	"testing"
)

// TestConfirmDestructive_Forceフラグが有効な場合は確認をスキップ
func TestConfirmDestructive_Forceフラグが有効な場合は確認をスキップ(t *testing.T) {
	// Forceフラグが有効な場合
	err := confirmDestructive("delete post 'テスト記事'", true, false)

	// エラーが発生しないことを検証
	if err != nil {
		t.Errorf("Force=true の場合はエラーが発生しない想定だが、エラーが返された: %v", err)
	}
}

// TestConfirmDestructive_Forceフラグが無効でも関数は存在する
func TestConfirmDestructive_Forceフラグが無効でも関数は存在する(t *testing.T) {
	// この時点では実装が不完全なため、関数が呼び出せることのみ確認
	// 実際の対話的入力テストは別途実装が必要
	_ = confirmDestructive("delete post 'テスト記事'", false, false)
}

// TestConfirmDestructive_NoInputフラグが有効な場合はForceなしでエラー
func TestConfirmDestructive_NoInputフラグが有効な場合はForceなしでエラー(t *testing.T) {
	// NoInputフラグが有効、Forceフラグが無効な場合
	err := confirmDestructive("delete post 'テスト記事'", false, true)

	// エラーが発生することを検証
	if err == nil {
		t.Error("NoInput=true かつ Force=false の場合はエラーが発生する想定だが、エラーが返されなかった")
	}

	// エラーメッセージに "non-interactive" が含まれることを検証
	if err != nil && !contains(err.Error(), "non-interactive") {
		t.Errorf("エラーメッセージに 'non-interactive' が含まれていない: %v", err)
	}
}

// TestConfirmDestructive_NoInputとForceの両方が有効な場合は確認をスキップ
func TestConfirmDestructive_NoInputとForceの両方が有効な場合は確認をスキップ(t *testing.T) {
	// NoInputフラグとForceフラグの両方が有効な場合
	err := confirmDestructive("delete post 'テスト記事'", true, true)

	// エラーが発生しないことを検証
	if err != nil {
		t.Errorf("Force=true の場合はエラーが発生しない想定だが、エラーが返された: %v", err)
	}
}

// contains は文字列に部分文字列が含まれているかチェックする補助関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
