/**
 * confirm.go
 * 破壊的操作の確認機構
 *
 * gogcliの安全機構を参考に、破壊的操作の実行前にユーザー確認を行う
 */

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/mtane0412/ghocli/internal/input"
	"github.com/mtane0412/ghocli/internal/ui"
)

// ConfirmDestructive は破壊的操作の実行前にユーザー確認を行います
//
// ctx: コンテキスト（UIの取得に使用）
// root: RootFlagsからForceとNoInputフラグを取得
// message: 実行する操作の説明（例: "delete post 'テスト記事'"）
//
// 戻り値:
//   - Force=true の場合は常にnilを返す
//   - NoInput=true かつ Force=false の場合はExitError{Code: 1}を返す
//   - 非対話的環境（TTYでない）の場合は、ForceなしではExitError{Code: 1}を返す
//   - 対話的環境では、ユーザーに確認プロンプトを表示し、y/yes以外の入力でExitError{Code: 1}を返す
func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	// Forceフラグが有効な場合は確認をスキップ
	if root.Force {
		return nil
	}

	// NoInputフラグが有効な場合、または非対話的環境の場合は、対話的入力を禁止
	if root.NoInput || !term.IsTerminal(int(os.Stdin.Fd())) {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("refusing to %s without --force (non-interactive)", message),
		}
	}

	// contextからUIを取得
	output := ui.FromContext(ctx)
	if output == nil {
		// UIが設定されていない場合はエラー（通常は発生しない）
		return &ExitError{
			Code: 1,
			Err:  errors.New("UI not configured in context"),
		}
	}

	// 確認プロンプトを表示し、ユーザー入力を読み取る
	prompt := fmt.Sprintf("Proceed to %s? [y/N]: ", message)
	line, readErr := input.PromptLineFrom(ctx, prompt, os.Stdin)
	if readErr != nil && !errors.Is(readErr, os.ErrClosed) {
		// EOFの場合はキャンセルとして扱う
		if errors.Is(readErr, io.EOF) {
			return &ExitError{Code: 1, Err: errors.New("cancelled")}
		}
		// その他のエラー
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("failed to read user input: %w", readErr),
		}
	}

	// 入力を正規化（前後の空白を削除、小文字化）
	ans := strings.ToLower(strings.TrimSpace(line))

	// "y" または "yes" の場合のみ続行
	if ans == "y" || ans == "yes" {
		return nil
	}

	// それ以外はキャンセル
	return &ExitError{Code: 1, Err: errors.New("cancelled")}
}
