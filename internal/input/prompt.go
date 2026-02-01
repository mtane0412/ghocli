/**
 * prompt.go
 * ユーザー入力プロンプト機能
 *
 * コマンドラインでユーザーに対話的に入力を求める機能を提供する。
 */
package input

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mtane0412/gho/internal/ui"
)

// PromptLine は、os.Stdinから1行の入力を読み取る
//
// ctx: コンテキスト（UIの取得に使用）
// prompt: ユーザーに表示するプロンプト文字列
//
// 戻り値:
//   - 読み取った行（改行文字を除く）
//   - エラー（読み取りに失敗した場合）
func PromptLine(ctx context.Context, prompt string) (string, error) {
	return PromptLineFrom(ctx, prompt, os.Stdin)
}

// PromptLineFrom は、指定されたio.Readerから1行の入力を読み取る
//
// ctx: コンテキスト（UIの取得に使用）
// prompt: ユーザーに表示するプロンプト文字列
// r: 入力元のio.Reader
//
// 戻り値:
//   - 読み取った行（改行文字を除く）
//   - エラー（読み取りに失敗した場合）
func PromptLineFrom(ctx context.Context, prompt string, r io.Reader) (string, error) {
	// contextからUIを取得し、プロンプトを出力
	if u := ui.FromContext(ctx); u != nil {
		// UIが設定されている場合は、PrintMessageを使用
		u.PrintMessage(prompt)
	} else {
		// UIが設定されていない場合は、stderrに直接出力
		_, _ = fmt.Fprint(os.Stderr, prompt)
	}

	// 1行読み取り
	return ReadLine(r)
}

// PromptPassword は、パスワード入力を促すプロンプトを表示する（TODO: 実装予定）
//
// この関数は将来的に実装予定。現時点では未実装。
func PromptPassword(ctx context.Context, prompt string) (string, error) {
	// TODO: terminal.ReadPasswordを使用してパスワード入力を実装
	return "", fmt.Errorf("PromptPassword is not implemented yet")
}
