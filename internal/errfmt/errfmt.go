/**
 * errfmtパッケージ
 *
 * エラーメッセージをユーザー向けにフォーマットする機能を提供する。
 * Ghost Admin API特有のエラーに対して、適切な対処法を含めたメッセージを生成する。
 */
package errfmt

import (
	"errors"
	"fmt"
)

// AuthRequiredError は、Ghost Admin APIへの認証が必要な場合のエラーを表す
type AuthRequiredError struct {
	// Site は認証が必要なGhostサイトのドメイン（例: "example.ghost.io"）
	Site string
	// Err は元のエラー
	Err error
}

// Error は、AuthRequiredErrorのエラーメッセージを返す
func (e *AuthRequiredError) Error() string {
	if e.Site != "" {
		return fmt.Sprintf("authentication required for %s", e.Site)
	}
	return "authentication required"
}

// Unwrap は、ラップされた元のエラーを返す
func (e *AuthRequiredError) Unwrap() error {
	return e.Err
}

// FormatAuthError は認証エラーをフォーマットする
//
// 認証が設定されていないサイトに対するエラーメッセージを生成し、
// gho auth addコマンドでの対処法を提示する。
func FormatAuthError(site string) string {
	return fmt.Sprintf(`No API key configured for site "%s".

Add credentials:
  gho auth add %s https://%s.ghost.io`, site, site, site)
}

// FormatSiteError はサイト未指定エラーをフォーマットする
//
// サイトが指定されていない場合のエラーメッセージを生成し、
// --siteフラグまたはdefault_site設定での対処法を提示する。
func FormatSiteError() string {
	return `No site specified.

Specify with --site flag or set default:
  gho config set default_site myblog`
}

// FormatFlagError は不明なフラグエラーをフォーマットする
//
// 不明なフラグが指定された場合のエラーメッセージを生成し、
// --helpフラグでの対処法を提示する。
func FormatFlagError(flag string) string {
	return fmt.Sprintf(`unknown flag %s
Run with --help to see available flags`, flag)
}

// Format は、エラーをユーザー向けにフォーマットする
//
// 以下の特別なエラー型を認識し、適切な対処法を含めたメッセージを返す：
// - AuthRequiredError: 認証エラー → gho auth loginコマンドを提示
// - その他のエラー: エラーメッセージをそのまま返す
//
// nilエラーの場合は空文字を返す。
func Format(err error) string {
	// nilエラーの場合は空文字を返す
	if err == nil {
		return ""
	}

	// AuthRequiredErrorの場合は、認証方法を提示する
	var authErr *AuthRequiredError
	if errors.As(err, &authErr) {
		if authErr.Site != "" {
			return fmt.Sprintf(
				"認証が必要です: %s\n\n対処法:\n  gho auth login %s",
				authErr.Site,
				authErr.Site,
			)
		}
		return "認証が必要です。\n\n対処法:\n  gho auth login <サイトURL>"
	}

	// その他のエラーはメッセージをそのまま返す
	return err.Error()
}
