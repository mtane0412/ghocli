/**
 * output.go
 * UI出力機能
 *
 * stdout/stderrの分離を提供します。
 * - データ出力 → stdout
 * - 進捗メッセージ → stderr
 * - エラーメッセージ → stderr
 */

package ui

import (
	"fmt"
	"io"
)

// Output はUI出力を管理する構造体です
type Output struct {
	stdout io.Writer
	stderr io.Writer
}

// NewOutput は新しいOutputを作成します
func NewOutput(stdout, stderr io.Writer) *Output {
	return &Output{
		stdout: stdout,
		stderr: stderr,
	}
}

// PrintData はデータをstdoutに出力します
func (o *Output) PrintData(data string) error {
	_, err := fmt.Fprintln(o.stdout, data)
	return err
}

// PrintMessage は進捗メッセージをstderrに出力します
func (o *Output) PrintMessage(message string) {
	fmt.Fprintln(o.stderr, message)
}

// PrintError はエラーメッセージをstderrに出力します
func (o *Output) PrintError(message string) {
	fmt.Fprintln(o.stderr, message)
}
