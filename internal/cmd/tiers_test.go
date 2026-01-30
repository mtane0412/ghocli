/**
 * tiers_test.go
 * ティア管理コマンドのテストコード
 */

package cmd

import (
	"testing"
)

// TestTiersInfoCmd_構造体が存在すること
func TestTiersInfoCmd_構造体が存在すること(t *testing.T) {
	// TiersInfoCmdが定義されていることを確認
	_ = &TiersInfoCmd{}
}
