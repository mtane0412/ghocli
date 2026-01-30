/**
 * color_test.go
 * カラー出力機能のテスト
 */

package ui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestShouldUseColor_NeverMode はNeverモードでカラー出力が無効になることをテストします
func TestShouldUseColor_NeverMode(t *testing.T) {
	// Neverモードでは常にfalseを返すべき
	assert.False(t, ShouldUseColor(ColorNever), "ColorNeverモードではfalseを返すべき")
}

// TestShouldUseColor_AlwaysMode はAlwaysモードでカラー出力が有効になることをテストします
func TestShouldUseColor_AlwaysMode(t *testing.T) {
	// NO_COLOR環境変数がない場合、Alwaysモードでtrueを返すべき
	os.Unsetenv("NO_COLOR")
	assert.True(t, ShouldUseColor(ColorAlways), "ColorAlwaysモードではtrueを返すべき")
}

// TestShouldUseColor_NO_COLOR はNO_COLOR環境変数が設定されている場合にカラー出力が無効になることをテストします
func TestShouldUseColor_NO_COLOR(t *testing.T) {
	// NO_COLOR環境変数を設定
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	// NO_COLORが設定されている場合、Alwaysモードでもfalseを返すべき
	assert.False(t, ShouldUseColor(ColorAlways), "NO_COLOR環境変数が設定されている場合、Alwaysモードでもfalseを返すべき")
}

// TestShouldUseColor_AutoMode はAutoモードでTTY判定が行われることをテストします
func TestShouldUseColor_AutoMode(t *testing.T) {
	// NO_COLOR環境変数がない場合
	os.Unsetenv("NO_COLOR")

	// Autoモードでは、TTY判定の結果を返す
	// CIではTTYではないので、通常はfalseになる
	// このテストでは、関数が呼び出せることのみ確認
	_ = ShouldUseColor(ColorAuto)
}
