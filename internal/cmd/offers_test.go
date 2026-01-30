/**
 * offers_test.go
 * オファー管理コマンドのテストコード
 *
 * Phase 2で追加される新規コマンドのテストを含みます。
 */

package cmd

import (
	"testing"
)

// TestOffersInfoCmd_構造体が存在すること
func TestOffersInfoCmd_構造体が存在すること(t *testing.T) {
	// OffersInfoCmdが定義されていることを確認
	_ = &OffersInfoCmd{}
}

// TestOffersArchiveCmd_構造体が存在すること
func TestOffersArchiveCmd_構造体が存在すること(t *testing.T) {
	// OffersArchiveCmdが定義されていることを確認
	_ = &OffersArchiveCmd{}
}
