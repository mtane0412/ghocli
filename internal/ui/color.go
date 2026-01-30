/**
 * color.go
 * カラー出力機能
 *
 * --colorフラグとNO_COLOR環境変数に対応したカラー出力制御を提供します。
 */

package ui

import (
	"os"

	"github.com/muesli/termenv"
)

// ColorMode はカラー出力のモードを表す型です
type ColorMode string

const (
	// ColorAuto はTTYの場合のみカラー出力を有効にします
	ColorAuto ColorMode = "auto"
	// ColorAlways は常にカラー出力を有効にします
	ColorAlways ColorMode = "always"
	// ColorNever はカラー出力を無効にします
	ColorNever ColorMode = "never"
)

// ShouldUseColor はカラー出力を使用すべきか判定します
func ShouldUseColor(mode ColorMode) bool {
	// Neverモードの場合は常にfalse
	if mode == ColorNever {
		return false
	}

	// NO_COLOR環境変数が設定されている場合はfalse
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Alwaysモードの場合は常にtrue
	if mode == ColorAlways {
		return true
	}

	// Autoモードの場合は、TTYかどうかで判定
	profile := termenv.NewOutput(os.Stdout).Profile
	return profile != termenv.Ascii
}
