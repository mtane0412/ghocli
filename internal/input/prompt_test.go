/**
 * prompt_test.go
 * ユーザー入力プロンプト機能のテストコード
 */
package input_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mtane0412/ghocli/internal/input"
	"github.com/mtane0412/ghocli/internal/ui"
)

// TestPromptLineFrom は、io.Readerから入力を読み取るPromptLineFrom関数のテスト
func TestPromptLineFrom(t *testing.T) {
	// 前提条件: stderrにプロンプトを出力するためのUIを準備
	var stderr bytes.Buffer
	output := ui.NewOutput(&stderr, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// 実行: PromptLineFrom関数を呼び出す（strings.Readerから読み取る）
	line, err := input.PromptLineFrom(ctx, "Prompt: ", strings.NewReader("hello\n"))

	// 検証: エラーが発生しないこと
	if err != nil {
		t.Fatalf("PromptLineFrom関数がエラーを返した: %v", err)
	}

	// 検証: 正しい行が読み取られること
	if line != "hello" {
		t.Errorf("読み取られた行が期待と異なる: got %q, want %q", line, "hello")
	}

	// 検証: プロンプトがstderrに出力されること
	if !strings.Contains(stderr.String(), "Prompt: ") {
		t.Errorf("プロンプトがstderrに出力されていない: %q", stderr.String())
	}
}

// TestPromptLine は、os.Stdinから入力を読み取るPromptLine関数のテスト
func TestPromptLine(t *testing.T) {
	// 前提条件: os.Stdinを保存し、テスト後に復元する
	orig := os.Stdin
	defer func() {
		os.Stdin = orig
	}()

	// 前提条件: パイプを作成してos.Stdinを置き換える
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("パイプの作成に失敗: %v", err)
	}
	defer func() {
		_ = r.Close()
	}()
	os.Stdin = r

	// 前提条件: パイプに入力データを書き込む
	_, writeErr := w.WriteString("world\n")
	if writeErr != nil {
		t.Fatalf("パイプへの書き込みに失敗: %v", writeErr)
	}
	_ = w.Close()

	// 実行: PromptLine関数を呼び出す
	line, err := input.PromptLine(context.Background(), "Prompt: ")

	// 検証: エラーが発生しないこと
	if err != nil {
		t.Fatalf("PromptLine関数がエラーを返した: %v", err)
	}

	// 検証: 正しい行が読み取られること
	if line != "world" {
		t.Errorf("読み取られた行が期待と異なる: got %q, want %q", line, "world")
	}
}
