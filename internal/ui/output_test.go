/**
 * output_test.go
 * UI出力機能のテスト
 */

package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewOutput_構造体が存在すること
func TestNewOutput_構造体が存在すること(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)
	assert.NotNil(t, output)
}

// TestOutput_PrintData はデータをstdoutに出力することをテストします
func TestOutput_PrintData(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	err := output.PrintData("test data")
	assert.NoError(t, err)

	// データはstdoutに出力される
	assert.Contains(t, stdout.String(), "test data")
	// stderrには何も出力されない
	assert.Empty(t, stderr.String())
}

// TestOutput_PrintMessage は進捗メッセージをstderrに出力することをテストします
func TestOutput_PrintMessage(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	output.PrintMessage("test message")

	// メッセージはstderrに出力される
	assert.Contains(t, stderr.String(), "test message")
	// stdoutには何も出力されない
	assert.Empty(t, stdout.String())
}

// TestOutput_PrintError はエラーメッセージをstderrに出力することをテストします
func TestOutput_PrintError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	output.PrintError("test error")

	// エラーメッセージはstderrに出力される
	assert.Contains(t, stderr.String(), "test error")
	// stdoutには何も出力されない
	assert.Empty(t, stdout.String())
}

// TestWithUI_contextにUIを埋め込む
func TestWithUI_contextにUIを埋め込む(t *testing.T) {
	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	ctx = WithUI(ctx, output)

	// contextからUIを取得できることを確認
	retrieved := FromContext(ctx)
	assert.NotNil(t, retrieved)
	assert.Equal(t, output, retrieved)
}

// TestFromContext_UIが設定されていない場合はnilを返す
func TestFromContext_UIが設定されていない場合はnilを返す(t *testing.T) {
	ctx := context.Background()

	retrieved := FromContext(ctx)
	assert.Nil(t, retrieved)
}

// TestFromContext_contextからUIを取得して出力できる
func TestFromContext_contextからUIを取得して出力できる(t *testing.T) {
	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	ctx = WithUI(ctx, output)

	// contextからUIを取得
	ui := FromContext(ctx)
	assert.NotNil(t, ui)

	// UIを使って出力
	err := ui.PrintData("test data from context")
	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), "test data from context")

	ui.PrintMessage("test message from context")
	assert.Contains(t, stderr.String(), "test message from context")
}
