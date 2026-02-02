/**
 * confirm_test.go
 * 破壊的操作の確認機構のテストコード
 */

package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/mtane0412/ghocli/internal/ui"
)

// TestConfirmDestructive_Forceフラグが有効な場合は確認をスキップ
func TestConfirmDestructive_Forceフラグが有効な場合は確認をスキップ(t *testing.T) {
	// 前提条件: Forceフラグが有効なRootFlagsを用意
	root := &RootFlags{
		Force:   true,
		NoInput: false,
	}

	// contextにUIを埋め込む
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// 実行: ConfirmDestructiveを呼び出す
	err := ConfirmDestructive(ctx, root, "delete post 'テスト記事'")

	// 検証: エラーが発生しないこと
	if err != nil {
		t.Errorf("Force=true の場合はエラーが発生しない想定だが、エラーが返された: %v", err)
	}
}

// TestConfirmDestructive_Forceフラグが無効でも関数は存在する
func TestConfirmDestructive_Forceフラグが無効でも関数は存在する(t *testing.T) {
	// 前提条件: Forceフラグが無効なRootFlagsを用意
	root := &RootFlags{
		Force:   false,
		NoInput: false,
	}

	// contextにUIを埋め込む
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// 実行: この時点では実装が不完全なため、関数が呼び出せることのみ確認
	// actualの対話的入力テストは別途実装が必要
	_ = ConfirmDestructive(ctx, root, "delete post 'テスト記事'")
}

// TestConfirmDestructive_NoInputフラグが有効な場合はForceなしでエラー
func TestConfirmDestructive_NoInputフラグが有効な場合はForceなしでエラー(t *testing.T) {
	// 前提条件: NoInputフラグが有効、Forceフラグが無効なRootFlagsを用意
	root := &RootFlags{
		Force:   false,
		NoInput: true,
	}

	// contextにUIを埋め込む
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// 実行: ConfirmDestructiveを呼び出す
	err := ConfirmDestructive(ctx, root, "delete post 'テスト記事'")

	// 検証: エラーが発生すること
	if err == nil {
		t.Error("NoInput=true かつ Force=false の場合はエラーが発生する想定だが、エラーが返されなかった")
	}

	// 検証: ExitErrorが返されること
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("ExitErrorが返される想定だが、異なるエラーが返された: %T %v", err, err)
	}

	// 検証: エラーメッセージに "non-interactive" が含まれること
	if err != nil && !contains(err.Error(), "non-interactive") {
		t.Errorf("エラーメッセージに 'non-interactive' が含まれていない: %v", err)
	}
}

// TestConfirmDestructive_NoInputとForceの両方が有効な場合は確認をスキップ
func TestConfirmDestructive_NoInputとForceの両方が有効な場合は確認をスキップ(t *testing.T) {
	// 前提条件: NoInputフラグとForceフラグの両方が有効なRootFlagsを用意
	root := &RootFlags{
		Force:   true,
		NoInput: true,
	}

	// contextにUIを埋め込む
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// 実行: ConfirmDestructiveを呼び出す
	err := ConfirmDestructive(ctx, root, "delete post 'テスト記事'")

	// 検証: エラーが発生しないこと
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
