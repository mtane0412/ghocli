/**
 * confirm.go
 * 破壊的操作の確認機構
 *
 * gogcliの安全機構を参考に、破壊的操作の実行前にユーザー確認を行う
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// confirmDestructive は破壊的操作の実行前にユーザー確認を行います
//
// action: 実行する操作の説明（例: "delete post 'テスト記事'"）
// force: 確認をスキップするフラグ
// noInput: 対話的入力を無効化するフラグ（CI環境向け）
//
// 戻り値:
//   - force=true の場合は常にnilを返す
//   - noInput=true かつ force=false の場合はエラーを返す
//   - 非対話的環境（TTYでない）の場合は、forceなしではエラーを返す
//   - 対話的環境では、ユーザーに確認プロンプトを表示し、y/yes以外の入力でエラーを返す
func confirmDestructive(action string, force bool, noInput bool) error {
	// Forceフラグが有効な場合は確認をスキップ
	if force {
		return nil
	}

	// NoInputフラグが有効な場合は、対話的入力を禁止
	if noInput {
		return fmt.Errorf("refusing to %s without --force (non-interactive)", action)
	}

	// TTY検出 - 非対話的環境では--forceなしで拒否
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("refusing to %s without --force (non-interactive)", action)
	}

	// 確認プロンプトを表示
	fmt.Printf("Proceed to %s? [y/N]: ", action)

	// ユーザー入力を読み取る
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	// 入力を正規化（前後の空白と改行を削除、小文字化）
	input = strings.ToLower(strings.TrimSpace(input))

	// "y" または "yes" 以外はキャンセル
	if input != "y" && input != "yes" {
		return fmt.Errorf("operation cancelled by user")
	}

	return nil
}
