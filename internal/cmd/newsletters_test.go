/**
 * newsletters_test.go
 * ニュースレター管理コマンドのテストコード
 */

package cmd

import (
	"testing"
)

// TestNewslettersInfoCmd_構造体が存在すること
func TestNewslettersInfoCmd_構造体が存在すること(t *testing.T) {
	// NewslettersInfoCmdが定義されていることを確認
	_ = &NewslettersInfoCmd{}
}
