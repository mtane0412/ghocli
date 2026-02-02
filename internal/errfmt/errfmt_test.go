/**
 * errfmtパッケージのテスト
 *
 * エラーフォーマット機能の動作を検証する。
 */
package errfmt_test

import (
	"errors"
	"testing"

	"github.com/mtane0412/ghocli/internal/errfmt"
)

// TestFormat_NilError は、nilエラーの場合に空文字を返すことを検証する
func TestFormat_NilError(t *testing.T) {
	// 実行
	result := errfmt.Format(nil)

	// 検証: 空文字が返されること
	if result != "" {
		t.Errorf("Format(nil) = %q; want empty string", result)
	}
}

// TestFormat_GenericError は、一般エラーの場合にそのままのメッセージを返すことを検証する
func TestFormat_GenericError(t *testing.T) {
	// 前提条件: 一般的なエラーを用意
	err := errors.New("something went wrong")

	// 実行
	result := errfmt.Format(err)

	// 検証: エラーメッセージがそのまま返されること
	expected := "something went wrong"
	if result != expected {
		t.Errorf("Format(generic error) = %q; want %q", result, expected)
	}
}

// TestFormat_AuthRequiredError は、認証エラーの場合に対処法付きメッセージを返すことを検証する
func TestFormat_AuthRequiredError(t *testing.T) {
	// 前提条件: 認証エラーを用意
	baseErr := errors.New("401 Unauthorized")
	authErr := &errfmt.AuthRequiredError{
		Site: "example.ghost.io",
		Err:  baseErr,
	}

	// 実行
	result := errfmt.Format(authErr)

	// 検証: 認証エラーメッセージに対処法が含まれること
	// メッセージには以下が含まれるべき：
	// - サイト名
	// - エラー内容
	// - 対処法（gho auth loginコマンド）
	if result == "" {
		t.Error("Format(AuthRequiredError) returned empty string")
	}

	// サイト名が含まれること
	if !contains(result, "example.ghost.io") {
		t.Errorf("Format(AuthRequiredError) = %q; want to contain site name 'example.ghost.io'", result)
	}

	// 対処法が含まれること
	if !contains(result, "gho auth login") {
		t.Errorf("Format(AuthRequiredError) = %q; want to contain 'gho auth login'", result)
	}
}

// contains は文字列sにsubstrが含まれるかを返すヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// TestFormatAuthError は認証エラーフォーマット関数をテストする
func TestFormatAuthError(t *testing.T) {
	// 実行
	msg := errfmt.FormatAuthError("myblog")

	// 検証: サイト名が含まれること
	if !contains(msg, "myblog") {
		t.Errorf("FormatAuthError() = %q; want to contain 'myblog'", msg)
	}

	// 検証: 認証エラーメッセージが含まれること
	if !contains(msg, "No API key configured") {
		t.Errorf("FormatAuthError() = %q; want to contain 'No API key configured'", msg)
	}

	// 検証: 対処法が含まれること
	if !contains(msg, "gho auth add myblog") {
		t.Errorf("FormatAuthError() = %q; want to contain 'gho auth add myblog'", msg)
	}
}

// TestFormatSiteError はサイト未指定エラーフォーマット関数をテストする
func TestFormatSiteError(t *testing.T) {
	// 実行
	msg := errfmt.FormatSiteError()

	// 検証: サイト未指定メッセージが含まれること
	if !contains(msg, "No site specified") {
		t.Errorf("FormatSiteError() = %q; want to contain 'No site specified'", msg)
	}

	// 検証: --siteフラグの説明が含まれること
	if !contains(msg, "--site") {
		t.Errorf("FormatSiteError() = %q; want to contain '--site'", msg)
	}

	// 検証: config set default_siteの説明が含まれること
	if !contains(msg, "gho config set default_site") {
		t.Errorf("FormatSiteError() = %q; want to contain 'gho config set default_site'", msg)
	}
}

// TestFormatFlagError は不明なフラグエラーフォーマット関数をテストする
func TestFormatFlagError(t *testing.T) {
	// 実行
	msg := errfmt.FormatFlagError("--foo")

	// 検証: 不明なフラグメッセージが含まれること
	if !contains(msg, "unknown flag --foo") {
		t.Errorf("FormatFlagError() = %q; want to contain 'unknown flag --foo'", msg)
	}

	// 検証: --helpヒントが含まれること
	if !contains(msg, "--help") {
		t.Errorf("FormatFlagError() = %q; want to contain '--help'", msg)
	}
}
