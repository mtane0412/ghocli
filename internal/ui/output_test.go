/**
 * output_test.go
 * UI出力機能のテスト
 */

package ui

import (
	"bytes"
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
